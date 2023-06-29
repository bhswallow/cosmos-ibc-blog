package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"planet/x/blog/types"
)

// GetUpdatePostCount get the total number of updatePost
func (k Keeper) GetUpdatePostCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.UpdatePostCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetUpdatePostCount set the total number of updatePost
func (k Keeper) SetUpdatePostCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.UpdatePostCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendUpdatePost appends a updatePost in the store with a new id and update the count
func (k Keeper) AppendUpdatePost(
	ctx sdk.Context,
	updatePost types.UpdatePost,
) uint64 {
	// Create the updatePost
	count := k.GetUpdatePostCount(ctx)

	// Set the ID of the appended value
	updatePost.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UpdatePostKey))
	appendedValue := k.cdc.MustMarshal(&updatePost)
	store.Set(GetUpdatePostIDBytes(updatePost.Id), appendedValue)

	// Update updatePost count
	k.SetUpdatePostCount(ctx, count+1)

	return count
}

// SetUpdatePost set a specific updatePost in the store
func (k Keeper) SetUpdatePost(ctx sdk.Context, updatePost types.UpdatePost) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UpdatePostKey))
	b := k.cdc.MustMarshal(&updatePost)
	store.Set(GetUpdatePostIDBytes(updatePost.Id), b)
}

// GetUpdatePost returns a updatePost from its id
func (k Keeper) GetUpdatePost(ctx sdk.Context, id uint64) (val types.UpdatePost, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UpdatePostKey))
	b := store.Get(GetUpdatePostIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUpdatePost removes a updatePost from the store
func (k Keeper) RemoveUpdatePost(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UpdatePostKey))
	store.Delete(GetUpdatePostIDBytes(id))
}

// GetAllUpdatePost returns all updatePost
func (k Keeper) GetAllUpdatePost(ctx sdk.Context) (list []types.UpdatePost) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UpdatePostKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UpdatePost
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetUpdatePostIDBytes returns the byte representation of the ID
func GetUpdatePostIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetUpdatePostIDFromBytes returns ID in uint64 format from a byte array
func GetUpdatePostIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
