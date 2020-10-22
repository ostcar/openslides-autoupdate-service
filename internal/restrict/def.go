// Code generated with models.txt DO NOT EDIT.
package restrict

var relationLists = map[string]string{
	"agenda_item/child_ids":                                    "agenda_item",
	"agenda_item/current_projector_ids":                        "projector",
	"agenda_item/projection_ids":                               "projection",
	"agenda_item/tag_ids":                                      "tag",
	"assignment/attachment_ids":                                "mediafile",
	"assignment/candidate_ids":                                 "assignment_candidate",
	"assignment/current_projector_ids":                         "projector",
	"assignment/poll_ids":                                      "assignment_poll",
	"assignment/projection_ids":                                "projection",
	"assignment/tag_ids":                                       "tag",
	"assignment_option/vote_ids":                               "assignment_vote",
	"assignment_poll/current_projector_ids":                    "projector",
	"assignment_poll/entitled_group_ids":                       "group",
	"assignment_poll/option_ids":                               "assignment_option",
	"assignment_poll/projection_ids":                           "projection",
	"assignment_poll/voted_ids":                                "user",
	"committee/forward_to_committee_ids":                       "committee",
	"committee/manager_ids":                                    "user",
	"committee/meeting_ids":                                    "meeting",
	"committee/member_ids":                                     "user",
	"committee/receive_forwardings_from_committee_ids":         "committee",
	"group/assignment_poll_ids":                                "assignment_poll",
	"group/mediafile_access_group_ids":                         "mediafile",
	"group/mediafile_inherited_access_group_ids":               "mediafile",
	"group/motion_poll_ids":                                    "motion_poll",
	"group/read_comment_section_ids":                           "motion_comment_section",
	"group/user_ids":                                           "user",
	"group/write_comment_section_ids":                          "motion_comment_section",
	"list_of_speakers/current_projector_ids":                   "projector",
	"list_of_speakers/projection_ids":                          "projection",
	"list_of_speakers/speaker_ids":                             "speaker",
	"mediafile/access_group_ids":                               "group",
	"mediafile/attachment_ids":                                 "*",
	"mediafile/child_ids":                                      "mediafile",
	"mediafile/current_projector_ids":                          "projector",
	"mediafile/inherited_access_group_ids":                     "group",
	"mediafile/projection_ids":                                 "projection",
	"mediafile/used_as_font_$_in_meeting_id":                   "meeting",
	"mediafile/used_as_logo_$_in_meeting_id":                   "meeting",
	"meeting/agenda_item_ids":                                  "agenda_item",
	"meeting/assignment_candidate_ids":                         "assignment_candidate",
	"meeting/assignment_ids":                                   "assignment",
	"meeting/assignment_option_ids":                            "assignment_option",
	"meeting/assignment_poll_default_group_ids":                "group",
	"meeting/assignment_poll_ids":                              "assignment_poll",
	"meeting/assignment_vote_ids":                              "assignment_vote",
	"meeting/font_$_id":                                        "mediafile",
	"meeting/group_ids":                                        "group",
	"meeting/guest_ids":                                        "user",
	"meeting/list_of_speakers_ids":                             "list_of_speakers",
	"meeting/logo_$_id":                                        "mediafile",
	"meeting/mediafile_ids":                                    "mediafile",
	"meeting/motion_block_ids":                                 "motion_block",
	"meeting/motion_category_ids":                              "motion_category",
	"meeting/motion_change_recommendation_ids":                 "motion_change_recommendation",
	"meeting/motion_comment_ids":                               "motion_comment",
	"meeting/motion_comment_section_ids":                       "motion_comment_section",
	"meeting/motion_ids":                                       "motion",
	"meeting/motion_option_ids":                                "motion_option",
	"meeting/motion_poll_default_group_ids":                    "group",
	"meeting/motion_poll_ids":                                  "motion_poll",
	"meeting/motion_state_ids":                                 "motion_state",
	"meeting/motion_statute_paragraph_ids":                     "motion_statute_paragraph",
	"meeting/motion_submitter_ids":                             "motion_submitter",
	"meeting/motion_vote_ids":                                  "motion_vote",
	"meeting/motion_workflow_ids":                              "motion_workflow",
	"meeting/personal_note_ids":                                "personal_note",
	"meeting/present_user_ids":                                 "user",
	"meeting/projection_ids":                                   "projection",
	"meeting/projectiondefault_ids":                            "projectiondefault",
	"meeting/projector_countdown_ids":                          "projector_countdown",
	"meeting/projector_ids":                                    "projector",
	"meeting/projector_message_ids":                            "projector_message",
	"meeting/speaker_ids":                                      "speaker",
	"meeting/tag_ids":                                          "tag",
	"meeting/temporary_user_ids":                               "user",
	"meeting/topic_ids":                                        "topic",
	"motion/amendment_ids":                                     "motion",
	"motion/attachment_ids":                                    "mediafile",
	"motion/change_recommendation_ids":                         "motion_change_recommendation",
	"motion/comment_ids":                                       "motion_comment",
	"motion/current_projector_ids":                             "projector",
	"motion/derived_motion_ids":                                "motion",
	"motion/personal_note_ids":                                 "personal_note",
	"motion/poll_ids":                                          "motion_poll",
	"motion/projection_ids":                                    "projection",
	"motion/recommendation_extension_reference_ids":            "motion",
	"motion/referenced_in_motion_recommendation_extension_ids": "motion",
	"motion/sort_child_ids":                                    "motion",
	"motion/submitter_ids":                                     "motion_submitter",
	"motion/supporter_ids":                                     "user",
	"motion/tag_ids":                                           "tag",
	"motion_block/current_projector_ids":                       "projector",
	"motion_block/motion_ids":                                  "motion",
	"motion_block/projection_ids":                              "projection",
	"motion_category/child_ids":                                "motion_category",
	"motion_category/motion_ids":                               "motion",
	"motion_comment_section/comment_ids":                       "motion_comment",
	"motion_comment_section/read_group_ids":                    "group",
	"motion_comment_section/write_group_ids":                   "group",
	"motion_option/vote_ids":                                   "motion_vote",
	"motion_poll/current_projector_ids":                        "projector",
	"motion_poll/entitled_group_ids":                           "group",
	"motion_poll/option_ids":                                   "motion_option",
	"motion_poll/projection_ids":                               "projection",
	"motion_poll/voted_ids":                                    "user",
	"motion_state/motion_ids":                                  "motion",
	"motion_state/motion_recommendation_ids":                   "motion",
	"motion_state/next_state_ids":                              "motion_state",
	"motion_state/previous_state_ids":                          "motion_state",
	"motion_statute_paragraph/motion_ids":                      "motion",
	"motion_workflow/state_ids":                                "motion_state",
	"organisation/committee_ids":                               "committee",
	"organisation/resource_ids":                                "resource",
	"organisation/role_ids":                                    "role",
	"projector/current_element_ids":                            "*",
	"projector/current_projection_ids":                         "projection",
	"projector/history_projection_ids":                         "projection",
	"projector/preview_projection_ids":                         "projection",
	"projector/projectiondefault_ids":                          "projectiondefault",
	"projector_countdown/current_projector_ids":                "projector",
	"projector_countdown/projection_ids":                       "projection",
	"projector_message/current_projector_ids":                  "projector",
	"projector_message/projection_ids":                         "projection",
	"role/user_ids":                                            "user",
	"tag/tagged_ids":                                           "*",
	"topic/attachment_ids":                                     "mediafile",
	"topic/current_projector_ids":                              "projector",
	"topic/projection_ids":                                     "projection",
	"topic/tag_ids":                                            "tag",
	"user/assignment_candidate_$_ids":                          "assignment_candidate",
	"user/assignment_delegated_vote_$_ids":                     "assignment_vote",
	"user/assignment_option_$_ids":                             "assignment_option",
	"user/assignment_poll_voted_$_ids":                         "assignment_poll",
	"user/assignment_vote_$_ids":                               "assignment_vote",
	"user/committee_as_manager_ids":                            "committee",
	"user/committee_as_member_ids":                             "committee",
	"user/current_projector_ids":                               "projector",
	"user/group_$_ids":                                         "group",
	"user/guest_meeting_ids":                                   "meeting",
	"user/is_present_in_meeting_ids":                           "meeting",
	"user/motion_delegated_vote_$_ids":                         "motion_vote",
	"user/motion_poll_voted_$_ids":                             "motion_poll",
	"user/motion_vote_$_ids":                                   "motion_vote",
	"user/personal_note_$_ids":                                 "personal_note",
	"user/projection_ids":                                      "projection",
	"user/speaker_$_ids":                                       "speaker",
	"user/submitted_motion_$_ids":                              "motion_submitter",
	"user/supported_motion_$_ids":                              "motion",
	"user/vote_delegated_$_to_id":                              "user",
	"user/vote_delegations_$_from_ids":                         "user",
}