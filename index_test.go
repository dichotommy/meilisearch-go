package meilisearch

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIndex_Delete(t *testing.T) {
	type args struct {
		createUid []string
		deleteUid []string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name:   "TestIndexDeleteOneIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:   "TestIndexDeleteOneIndexWithCustomClient",
			client: customClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:   "TestIndexDeleteMultipleIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"1", "2", "3", "4", "5"},
				deleteUid: []string{"1", "2", "3", "4", "5"},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:   "TestIndexDeleteNotExistingIndex",
			client: defaultClient,
			args: args{
				createUid: []string{},
				deleteUid: []string{"1"},
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name:   "TestIndexDeleteMultipleNotExistingIndex",
			client: defaultClient,
			args: args{
				createUid: []string{},
				deleteUid: []string{"1", "2", "3"},
			},
			wantOk:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			for _, uid := range tt.args.createUid {
				_, err := c.CreateIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in DeleteTest error should be nil")
			}
			for _, uid := range tt.args.deleteUid {
				i := c.Index(uid)
				gotOk, err := i.Delete(uid)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.Equal(t, tt.wantOk, gotOk)
				}
			}
		})
	}
}

func TestIndex_GetStats(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *StatsIndex
	}{
		{
			name: "TestIndexBasicGetStats",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &StatsIndex{
				NumberOfDocuments: 6,
				IsIndexing:        false,
			},
		},
		{
			name: "TestIndexGetStatsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &StatsIndex{
				NumberOfDocuments: 6,
				IsIndexing:        false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpBasicIndex()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetStats()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			deleteAllIndexes(c)
		})
	}
}

func Test_newIndex(t *testing.T) {
	type args struct {
		client *Client
		uid    string
	}
	tests := []struct {
		name string
		args args
		want *Index
	}{
		{
			name: "TestBasicNewIndex",
			args: args{
				client: defaultClient,
				uid:    "TestBasicNewIndex",
			},
			want: &Index{
				UID:    "TestBasicNewIndex",
				client: defaultClient,
			},
		},
		{
			name: "TestNewIndexCustomClient",
			args: args{
				client: customClient,
				uid:    "TestBasicNewIndex",
			},
			want: &Index{
				UID:    "TestBasicNewIndex",
				client: customClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client

			got := newIndex(c, tt.args.uid)
			require.Equal(t, tt.want.UID, got.UID)
			require.Equal(t, tt.want.client, got.client)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetUpdateStatus(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		updateID int64
		document []docTest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetUpdateStatus",
			args: args{
				UID:      "1",
				client:   defaultClient,
				updateID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetUpdateStatusWithCustomClient",
			args: args{
				UID:      "1",
				client:   customClient,
				updateID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetUpdateStatus",
			args: args{
				UID:      "1",
				client:   defaultClient,
				updateID: 1,
				document: []docTest{
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			update, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			gotResp, err := i.GetUpdateStatus(update.UpdateID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.GreaterOrEqual(t, gotResp.UpdateID, tt.args.updateID)
			require.NotNil(t, gotResp.UpdateID)
		})
	}
	_, _ = deleteAllIndexes(defaultClient)
}

func TestIndex_GetAllUpdateStatus(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		document []docTest
	}
	tests := []struct {
		name     string
		args     args
		wantResp []Update
	}{
		{
			name: "TestIndexBasicGetAllUpdateStatus",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: []Update{
				{
					Status:   "processed",
					UpdateID: 0,
					Error:    "",
				},
			},
		},
		{
			name: "TestIndexGetAllUpdateStatusWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: []Update{
				{
					Status:   "processed",
					UpdateID: 0,
					Error:    "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			SetUpBasicIndex()

			gotResp, err := i.GetAllUpdateStatus()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp[0].Status, (*gotResp)[0].Status)
			require.Equal(t, tt.wantResp[0].UpdateID, (*gotResp)[0].UpdateID)
			require.Equal(t, tt.wantResp[0].Error, (*gotResp)[0].Error)
		})
	}
}

func TestIndex_DefaultWaitForPendingUpdate(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		updateID *AsyncUpdateID
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want UpdateStatus
	}{
		{
			name: "TestDefaultWaitForPendingUpdate",
			args: args{
				UID:    "1",
				client: defaultClient,
				updateID: &AsyncUpdateID{
					UpdateID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			want: "processed",
		},
		{
			name: "TestDefaultWaitForPendingUpdateWithCustomClient",
			args: args{
				UID:    "1",
				client: customClient,
				updateID: &AsyncUpdateID{
					UpdateID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			want: "processed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			update, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			got, err := i.DefaultWaitForPendingUpdate(update)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_WaitForPendingUpdate(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		interval time.Duration
		timeout  time.Duration
		updateID *AsyncUpdateID
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want UpdateStatus
	}{
		{
			name: "TestDefaultWaitForPendingUpdate50",
			args: args{
				UID:      "1",
				client:   defaultClient,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				updateID: &AsyncUpdateID{
					UpdateID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "processed",
		},
		{
			name: "TestDefaultWaitForPendingUpdate50WithCustomClient",
			args: args{
				UID:      "1",
				client:   customClient,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				updateID: &AsyncUpdateID{
					UpdateID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "processed",
		},
		{
			name: "TestDefaultWaitForPendingUpdate10",
			args: args{
				UID:      "1",
				client:   defaultClient,
				interval: time.Millisecond * 10,
				timeout:  time.Second * 5,
				updateID: &AsyncUpdateID{
					UpdateID: 1,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "processed",
		},
		{
			name: "TestDefaultWaitForPendingUpdateWithTimeout",
			args: args{
				UID:      "1",
				client:   defaultClient,
				interval: time.Millisecond * 50,
				timeout:  time.Millisecond * 10,
				updateID: &AsyncUpdateID{
					UpdateID: 1,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "processed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			update, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			ctx, cancelFunc := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancelFunc()

			got, err := i.WaitForPendingUpdate(ctx, tt.args.interval, update)
			if tt.args.timeout < tt.args.interval {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_FetchInfo(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Index
	}{
		{
			name: "TestIndexBasicFetchInfo",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &Index{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexFetchInfoWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &Index{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpBasicIndex()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.FetchInfo()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp.UID, gotResp.UID)
			require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_FetchPrimaryKey(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name           string
		args           args
		wantPrimaryKey string
	}{
		{
			name: "TestIndexBasicFetchPrimaryKey",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantPrimaryKey: "book_id",
		},
		{
			name: "TestIndexFetchPrimaryKeyWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantPrimaryKey: "book_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpBasicIndex()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotPrimaryKey, err := i.FetchPrimaryKey()
			require.NoError(t, err)
			require.Equal(t, &tt.wantPrimaryKey, gotPrimaryKey)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateIndex(t *testing.T) {
	type args struct {
		primaryKey string
		config     IndexConfig
		client     *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Index
	}{
		{
			name: "TestIndexBasicUpdateIndex",
			args: args{
				client: defaultClient,
				config: IndexConfig{
					Uid: "indexUID",
				},
				primaryKey: "book_id",
			},
			wantResp: &Index{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexUpdateIndexWithCustomClient",
			args: args{
				client: customClient,
				config: IndexConfig{
					Uid: "indexUID",
				},
				primaryKey: "book_id",
			},
			wantResp: &Index{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i, err := c.CreateIndex(&tt.args.config)
			require.NoError(t, err)
			require.Equal(t, tt.args.config.Uid, i.UID)
			require.Equal(t, tt.args.config.PrimaryKey, i.PrimaryKey)

			gotResp, err := i.UpdateIndex(tt.args.primaryKey)

			require.NoError(t, err)
			require.Equal(t, tt.wantResp.UID, gotResp.UID)
			require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)

			deleteAllIndexes(c)
		})
	}
}