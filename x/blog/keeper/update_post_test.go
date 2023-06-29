package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "planet/testutil/keeper"
	"planet/testutil/nullify"
	"planet/x/blog/keeper"
	"planet/x/blog/types"
)

func createNUpdatePost(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.UpdatePost {
	items := make([]types.UpdatePost, n)
	for i := range items {
		items[i].Id = keeper.AppendUpdatePost(ctx, items[i])
	}
	return items
}

func TestUpdatePostGet(t *testing.T) {
	keeper, ctx := keepertest.BlogKeeper(t)
	items := createNUpdatePost(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetUpdatePost(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestUpdatePostRemove(t *testing.T) {
	keeper, ctx := keepertest.BlogKeeper(t)
	items := createNUpdatePost(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveUpdatePost(ctx, item.Id)
		_, found := keeper.GetUpdatePost(ctx, item.Id)
		require.False(t, found)
	}
}

func TestUpdatePostGetAll(t *testing.T) {
	keeper, ctx := keepertest.BlogKeeper(t)
	items := createNUpdatePost(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllUpdatePost(ctx)),
	)
}

func TestUpdatePostCount(t *testing.T) {
	keeper, ctx := keepertest.BlogKeeper(t)
	items := createNUpdatePost(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetUpdatePostCount(ctx))
}
