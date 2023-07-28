package collection

import (
	"context"
	"fmt"

	"github.com/OpenSlides/openslides-autoupdate-service/internal/restrict/perm"
	"github.com/OpenSlides/openslides-autoupdate-service/internal/restrict2/attribute"
	"github.com/OpenSlides/openslides-autoupdate-service/pkg/datastore/dsfetch"
)

// User handels the restrictions for the user collection.
//
// Y is the request user and X the user, that is requested.
//
// The user Y can see a user X, if at least one condition is true:
//
//	Y==X
//	Y has the OML can_manage_users or higher.
//	There exists a committee where Y has the CML can_manage and X is in committee/user_ids.
//	X is in a group of a meeting where Y has user.can_see.
//	There exists a meeting where Y has the CML can_manage for the meeting's committee X is in meeting/user_ids.
//	There is a related object:
//	    There exists a motion which Y can see and X is a submitter/supporter.
//	    There exists an option which Y can see and X is the linked content object.
//	    There exists an assignment candidate which Y can see and X is the linked user.
//	    There exists a speaker which Y can see and X is the linked user.
//	    There exists a poll where Y can see the poll/voted_ids and X is part of that list.
//	    There exists a vote which Y can see and X is linked in user_id or delegated_user_id.
//	    There exists a chat_message which Y can see and X has sent it (specified by chat_message/user_id).
//	X is linked in one of the relations vote_delegated_$_to_id or vote_delegations_$_from_ids of Y.
//
// Mode A: Y can see X.
//
// Mode B: Y==X.
//
// Mode D: Y can see these fields if at least one condition is true:
//
//	Y has the OML can_manage_users or higher.
//	X is in a group of a meeting where Y has user.can_manage.
//
// Mode E: Y can see these fields if at least one condition is true:
//
//	Y has the OML can_manage_users or higher.
//	There exists a committee where Y has the CML can_manage and X is in committee/user_ids.
//	X is in a group of a meeting where Y has user.can_manage.
//	Y==X.
//
// Mode F: Y has the OML can_manage_users or higher.
//
// Mode G: No one. Not even the superadmin.
//
// Mode H: Like D but the fields are not visible, if the request has a lower
// organization management level then the requested user.
type User struct{}

// Name returns the collection name.
func (u User) Name() string {
	return "user"
}

// MeetingID returns the meetingID for the object.
func (User) MeetingID(ctx context.Context, ds *dsfetch.Fetch, id int) (int, bool, error) {
	return 0, false, nil
}

// Modes returns the field restriction for each mode.
func (u User) Modes(mode string) FieldRestricter {
	switch mode {
	case "A":
		return u.see
	case "B":
		return u.modeB
	case "D":
		return u.modeD
	case "E":
		return u.modeE
	case "F":
		return u.modeF
	case "G":
		return never
	case "H":
		return u.modeH
	}
	return nil
}

// SuperAdmin restricts the super admin.
func (u User) SuperAdmin(mode string) FieldRestricter {
	// TODO: When to call me????
	if mode == "G" {
		return never
	}
	return Allways
}

// TODO: this is not good.
func (u User) see(ctx context.Context, fetcher *dsfetch.Fetch, userIDs []int) ([]attribute.Func, error) {
	userManager := attribute.FuncGlobalLevel(perm.OMLCanManageUsers)

	inCommitteList := make([][]int, len(userIDs))
	inMeetingList := make([][]int, len(userIDs))
	for i, userID := range userIDs {
		fetcher.User_CommitteeIDs(userID).Lazy(&inCommitteList[i])
		fetcher.User_MeetingIDs(userID).Lazy(&inMeetingList[i])
	}

	if err := fetcher.Execute(ctx); err != nil {
		return nil, fmt.Errorf("fetching user data: %w", err)
	}

	result := make([]attribute.Func, len(userIDs))
	for i, userID := range userIDs {
		committeeManagerList := make([][]int, len(inCommitteList[i]))
		for j, committeeID := range inCommitteList[i] {
			fetcher.Committee_UserManagementLevel(committeeID, "can_manage").Lazy(&committeeManagerList[j])
		}

		if err := fetcher.Execute(ctx); err != nil {
			return nil, fmt.Errorf("fetching committee manager for all committees of user %d: %w", userID, err)
		}

		var committeeManagers []int
		for _, l := range committeeManagerList {
			committeeManagers = append(committeeManagers, l...)
		}

		var canSeeGroups []int
		for _, meetingID := range inMeetingList[i] {
			groupMap, err := perm.GroupMapFromContext(ctx, fetcher, meetingID)
			if err != nil {
				return nil, fmt.Errorf("group map for user %d in meeting %d: %w", userID, meetingID, err)
			}

			canSeeGroups = append(canSeeGroups, groupMap[perm.UserCanSee]...)
		}

		result[i] = attribute.FuncOr(
			attribute.FuncUserIDs([]int{userID}),
			userManager,
			attribute.FuncUserIDs(committeeManagers),
			attribute.FuncInGroup(canSeeGroups),
			// TODO There exists a meeting where Y has the CML can_manage for the meeting's committee X is in meeting/user_ids.
			// ...
		)
	}
	return result, nil
}

// UserRequiredObject represents the reference from a user to other objects.
type UserRequiredObject struct {
	Name     string
	TmplFunc func(int) *dsfetch.ValueIDSlice
	ElemFunc func(int, int) *dsfetch.ValueIntSlice
	SeeFunc  FieldRestricter
}

// RequiredObjects returns all references to other objects from the user.
func (User) RequiredObjects(ctx context.Context, ds *dsfetch.Fetch) []UserRequiredObject {
	return []UserRequiredObject{
		{
			"motion submitter",
			ds.User_SubmittedMotionIDsTmpl,
			ds.User_SubmittedMotionIDs,
			Collection(ctx, "motion_submitter").Modes("A"),
		},

		{
			"motion supporter",
			ds.User_SupportedMotionIDsTmpl,
			ds.User_SupportedMotionIDs,
			Collection(ctx, Motion{}.Name()).Modes("C"),
		},

		// {
		// 	"option",
		// 	ds.User_OptionIDsTmpl,
		// 	ds.User_OptionIDs,
		// 	Collection(ctx, Option{}.Name()).Modes("A"),
		// },

		// {
		// 	"assignment candidate",
		// 	ds.User_AssignmentCandidateIDsTmpl,
		// 	ds.User_AssignmentCandidateIDs,
		// 	Collection(ctx, AssignmentCandidate{}.Name()).Modes("A"),
		// },

		// {
		// 	"speaker",
		// 	ds.User_SpeakerIDsTmpl,
		// 	ds.User_SpeakerIDs,
		// 	Collection(ctx, Speaker{}.Name()).Modes("A"),
		// },

		// {
		// 	"poll voted",
		// 	ds.User_PollVotedIDsTmpl,
		// 	ds.User_PollVotedIDs,
		// 	Collection(ctx, Poll{}.Name()).Modes("A"),
		// },

		// {
		// 	"vote user",
		// 	ds.User_VoteIDsTmpl,
		// 	ds.User_VoteIDs,
		// 	Collection(ctx, Vote{}.Name()).Modes("A"),
		// },

		// {
		// 	"vote delegated user",
		// 	ds.User_VoteDelegatedVoteIDsTmpl,
		// 	ds.User_VoteDelegatedVoteIDs,
		// 	Collection(ctx, Vote{}.Name()).Modes("A"),
		// },

		// {
		// 	"chat messages",
		// 	ds.User_ChatMessageIDsTmpl,
		// 	ds.User_ChatMessageIDs,
		// 	Collection(ctx, ChatMessage{}.Name()).Modes("A"),
		// },
	}
}

func (u User) modeB(ctx context.Context, fetcher *dsfetch.Fetch, userIDs []int) ([]attribute.Func, error) {
	result := make([]attribute.Func, len(userIDs))
	for i, userID := range userIDs {
		result[i] = attribute.FuncUserIDs([]int{userID})
	}
	return result, nil
}

func (u User) modeD(ctx context.Context, fetcher *dsfetch.Fetch, userIDs []int) ([]attribute.Func, error) {
	userManager := attribute.FuncGlobalLevel(perm.OMLCanManageUsers)

	inMeetingList := make([][]int, len(userIDs))
	for i, userID := range userIDs {
		fetcher.User_MeetingIDs(userID).Lazy(&inMeetingList[i])
	}

	if err := fetcher.Execute(ctx); err != nil {
		return nil, fmt.Errorf("fetching user data: %w", err)
	}

	result := make([]attribute.Func, len(userIDs))
	for i, userID := range userIDs {
		var canSeeGroups []int
		for _, meetingID := range inMeetingList[i] {
			groupMap, err := perm.GroupMapFromContext(ctx, fetcher, meetingID)
			if err != nil {
				return nil, fmt.Errorf("group map for user %d in meeting %d: %w", userID, meetingID, err)
			}

			canSeeGroups = append(canSeeGroups, groupMap[perm.UserCanManage]...)
		}

		result[i] = attribute.FuncOr(
			userManager,
			attribute.FuncInGroup(canSeeGroups),
		)
	}
	return result, nil
}

func (u User) modeE(ctx context.Context, fetcher *dsfetch.Fetch, userIDs []int) ([]attribute.Func, error) {
	userManager := attribute.FuncGlobalLevel(perm.OMLCanManageUsers)

	inCommitteList := make([][]int, len(userIDs))
	inMeetingList := make([][]int, len(userIDs))
	for i, userID := range userIDs {
		fetcher.User_CommitteeIDs(userID).Lazy(&inCommitteList[i])
		fetcher.User_MeetingIDs(userID).Lazy(&inMeetingList[i])
	}

	if err := fetcher.Execute(ctx); err != nil {
		return nil, fmt.Errorf("fetching user data: %w", err)
	}

	result := make([]attribute.Func, len(userIDs))
	for i, userID := range userIDs {
		committeeManagerList := make([][]int, len(inCommitteList[i]))
		for j, committeeID := range inCommitteList[i] {
			fetcher.Committee_UserManagementLevel(committeeID, "can_manage").Lazy(&committeeManagerList[j])
		}

		if err := fetcher.Execute(ctx); err != nil {
			return nil, fmt.Errorf("fetching committee manager for all committees of user %d: %w", userID, err)
		}

		var committeeManagers []int
		for _, l := range committeeManagerList {
			committeeManagers = append(committeeManagers, l...)
		}

		var canSeeGroups []int
		for _, meetingID := range inMeetingList[i] {
			groupMap, err := perm.GroupMapFromContext(ctx, fetcher, meetingID)
			if err != nil {
				return nil, fmt.Errorf("group map for user %d in meeting %d: %w", userID, meetingID, err)
			}

			canSeeGroups = append(canSeeGroups, groupMap[perm.UserCanSee]...)
		}

		result[i] = attribute.FuncOr(
			attribute.FuncUserIDs([]int{userID}),
			userManager,
			attribute.FuncUserIDs(committeeManagers),
			attribute.FuncInGroup(canSeeGroups),
		)
	}
	return result, nil
}

func (u User) modeF(ctx context.Context, fetcher *dsfetch.Fetch, userIDs []int) ([]attribute.Func, error) {
	userManager := attribute.FuncGlobalLevel(perm.OMLCanManageUsers)
	return attributeFuncList(len(userIDs), userManager), nil
}

// higherOrgaManagement returns true if request equal or higher  then
// request.
//
// An empty string is a valid organization management level for this function
// that has the lowest value.
func higherOrgaManagement(level perm.OrganizationManagementLevel) perm.OrganizationManagementLevel {
	switch level {
	case perm.OMLNone:
		return perm.OMLCanManageUsers
	case perm.OMLCanManageUsers:
		return perm.OMLCanManageOrganization
	default:
		return perm.OMLSuperadmin
	}
}

func (u User) modeH(ctx context.Context, fetcher *dsfetch.Fetch, userIDs []int) ([]attribute.Func, error) {
	userOrgaLevel := make([]string, len(userIDs))
	for i, userID := range userIDs {
		fetcher.User_OrganizationManagementLevel(userID).Lazy(&userOrgaLevel[i])
	}

	if err := fetcher.Execute(ctx); err != nil {
		return nil, fmt.Errorf("fetching orga levels: %w", err)
	}

	result, err := Collection(ctx, "user").Modes("D")(ctx, fetcher, userIDs)
	if err != nil {
		return nil, fmt.Errorf("check like D: %w", err)
	}

	for i := range userIDs {
		result[i] = attribute.FuncAnd(
			attribute.FuncGlobalLevel(higherOrgaManagement(perm.OrganizationManagementLevel(userOrgaLevel[i]))),
			result[i],
		)
	}
	return result, nil
}