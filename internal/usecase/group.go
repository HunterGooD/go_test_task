package usecase

import (
	"context"

	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
)

type GroupUsecase struct {
	groupRepo                     interfaces.GroupRepository
	transactionManagerSongsGroups interfaces.TransactionManagerSongsGroups
	log                           interfaces.Logger
}

// NewGroupUsecase operations with songs get, change, delete
func NewGroupUsecase(gr interfaces.GroupRepository, tmSG interfaces.TransactionManagerSongsGroups, logger interfaces.Logger) *GroupUsecase {
	logger.Info("create group usecase")
	return &GroupUsecase{gr, tmSG, logger}
}

func (gu *GroupUsecase) DeleteSoftGroup(ctx context.Context, group_name string) error {
	// TODO: DO implement
	return nil
	err := gu.transactionManagerSongsGroups.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			gu.transactionManagerSongsGroups.Rollback()
		}
	}()
	// -- SOFT DELETE groups

	// SELECT id
	// 	FROM groups
	// 	WHERE g_name='rockers';

	// UPDATE groups
	// 	SET deleted_at=CURRENT_TIMESTAMP
	// 	WHERE id=1;

	// UPDATE songs
	// 	SET deleted_at=CURRENT_TIMESTAMP
	// 	WHERE group_id=1;

	groupTX := gu.transactionManagerSongsGroups.GroupRepository()
	songTX := gu.transactionManagerSongsGroups.SongRepository()

	group, err := groupTX.GetByName(ctx, group_name)
	if err != nil {
		// gu.transactionManagerSongsGroups.Rollback()
		return err
	}

	if err = groupTX.DeleteSoftByID(ctx, group.ID); err != nil {
		// gu.transactionManagerSongsGroups.Rollback()
		return err
	}

	if err = songTX.DeleteSoftByGroupID(ctx, group.ID); err != nil {
		// gu.transactionManagerSongsGroups.Rollback()
		return err
	}

	return gu.transactionManagerSongsGroups.Commit()
}
