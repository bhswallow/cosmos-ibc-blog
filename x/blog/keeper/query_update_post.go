package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"planet/x/blog/types"
)

func (k Keeper) UpdatePostAll(goCtx context.Context, req *types.QueryAllUpdatePostRequest) (*types.QueryAllUpdatePostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var updatePosts []types.UpdatePost
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	updatePostStore := prefix.NewStore(store, types.KeyPrefix(types.UpdatePostKey))

	pageRes, err := query.Paginate(updatePostStore, req.Pagination, func(key []byte, value []byte) error {
		var updatePost types.UpdatePost
		if err := k.cdc.Unmarshal(value, &updatePost); err != nil {
			return err
		}

		updatePosts = append(updatePosts, updatePost)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllUpdatePostResponse{UpdatePost: updatePosts, Pagination: pageRes}, nil
}

func (k Keeper) UpdatePost(goCtx context.Context, req *types.QueryGetUpdatePostRequest) (*types.QueryGetUpdatePostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	updatePost, found := k.GetUpdatePost(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetUpdatePostResponse{UpdatePost: updatePost}, nil
}
