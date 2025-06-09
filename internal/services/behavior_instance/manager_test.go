package behavior_instance

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/common"
	"testing"
)

var response = `{
    "errors": [],
    "data": [
        {
            "id": "47c33f80-f2ca-4d2a-a73f-6b9195e68f45",
            "type": "behavior_instance",
            "script_id": "bba79770-19f1-11ec-9621-0242ac130002",
            "enabled": true,
            "state": {
                "source_type": "device",
                "model_id": "SML003"
            },
            "configuration": {
                "settings": {
                    "daylight_sensitivity": {
                        "dark_threshold": 14833,
                        "offset": 7000
                    }
                },
                "source": {
                    "rid": "f1fcd1cb-b439-4668-b6a9-0b57127af314",
                    "rtype": "device"
                },
                "when": {
                    "timeslots": [
                        {
                            "on_motion": {
                                "recall_single": [
                                    {
                                        "action": {
                                            "recall": {
                                                "rid": "68b39f81-1c15-4c82-bd0b-ab28606f3d2e",
                                                "rtype": "scene"
                                            }
                                        }
                                    }
                                ]
                            },
                            "on_no_motion": {
                                "after": {
                                    "minutes": 10
                                },
                                "recall_single": [
                                    {
                                        "action": "all_off"
                                    }
                                ]
                            },
                            "start_time": {
                                "time": {
                                    "hour": 8,
                                    "minute": 0
                                },
                                "type": "time"
                            }
                        },
                        {
                            "on_motion": {
                                "recall_single": [
                                    {
                                        "action": {
                                            "recall": {
                                                "rid": "68b39f81-1c15-4c82-bd0b-ab28606f3d2e",
                                                "rtype": "scene"
                                            }
                                        }
                                    }
                                ]
                            },
                            "on_no_motion": {
                                "after": {
                                    "minutes": 5
                                },
                                "recall_single": [
                                    {
                                        "action": "all_off"
                                    }
                                ]
                            },
                            "start_time": {
                                "time": {
                                    "hour": 23,
                                    "minute": 0
                                },
                                "type": "time"
                            }
                        }
                    ]
                },
                "where": [
                    {
                        "group": {
                            "rid": "0d960eab-68c6-4ed7-8c0d-a24ca756d58e",
                            "rtype": "room"
                        }
                    }
                ]
            },
            "dependees": [
                {
                    "target": {
                        "rid": "f1fcd1cb-b439-4668-b6a9-0b57127af314",
                        "rtype": "device"
                    },
                    "level": "critical",
                    "type": "ResourceDependee"
                },
                {
                    "target": {
                        "rid": "0d960eab-68c6-4ed7-8c0d-a24ca756d58e",
                        "rtype": "room"
                    },
                    "level": "critical",
                    "type": "ResourceDependee"
                },
                {
                    "target": {
                        "rid": "68b39f81-1c15-4c82-bd0b-ab28606f3d2e",
                        "rtype": "scene"
                    },
                    "level": "critical",
                    "type": "ResourceDependee"
                },
                {
                    "target": {
                        "rid": "961975da-d821-4165-9bba-6f5495a9feef",
                        "rtype": "temperature"
                    },
                    "level": "critical",
                    "type": "ResourceDependee"
                },
                {
                    "target": {
                        "rid": "deb7e1cd-2bb8-46a8-a5c2-0f2ec387c235",
                        "rtype": "light_level"
                    },
                    "level": "critical",
                    "type": "ResourceDependee"
                },
                {
                    "target": {
                        "rid": "0ee6be79-9078-4352-ae1b-ce6d31f4cb8c",
                        "rtype": "motion"
                    },
                    "level": "critical",
                    "type": "ResourceDependee"
                }
            ],
            "status": "running",
            "last_error": "",
            "metadata": {
                "name": "Bathroom"
            }
        }
    ]
}
`

func TestManager_JSONParse(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		expected := behavior_instance.Data{
			ID:       "47c33f80-f2ca-4d2a-a73f-6b9195e68f45",
			Type:     "behavior_instance",
			ScriptID: "bba79770-19f1-11ec-9621-0242ac130002",
			Enabled:  true,
			Status:   "running",
			Metadata: behavior_instance.Metadata{
				Name: "Bathroom",
			},
			State: behavior_instance.State{
				SourceType: "device",
				ModelID:    "SML003",
			},
			Configuration: behavior_instance.Configuration{
				Settings: behavior_instance.Settings{
					DaylightSensitivity: behavior_instance.DaylightSensitivity{
						DarkThreshold: 14833,
						Offset:        7000,
					},
				},
				Where: []behavior_instance.Where{
					{
						Group: common.Reference{
							RID:   "0d960eab-68c6-4ed7-8c0d-a24ca756d58e",
							RType: "room",
						},
					},
				},
				When: behavior_instance.When{
					Timeslots: []behavior_instance.TimeSlots{
						{
							OnMotion: behavior_instance.OnMotion{
								RecallSingle: []behavior_instance.RecallSingle{
									{
										Action: behavior_instance.Action{
											Recall: common.Reference{
												RID:   "68b39f81-1c15-4c82-bd0b-ab28606f3d2e",
												RType: "scene",
											},
										},
									},
								},
							},
							OnNoMotion: behavior_instance.OnNoMotion{
								After: behavior_instance.After{
									Minutes: 10,
								},
								RecallSingle: []behavior_instance.RecallSingleNoMotion{
									{
										Action: "all_off",
									},
								},
							},
							StartTime: behavior_instance.StartTime{
								Time: behavior_instance.Time{
									Hour:   8,
									Minute: 0,
								},
								Type: "time",
							},
						},
						{
							OnMotion: behavior_instance.OnMotion{
								RecallSingle: []behavior_instance.RecallSingle{
									{
										Action: behavior_instance.Action{
											Recall: common.Reference{
												RID:   "68b39f81-1c15-4c82-bd0b-ab28606f3d2e",
												RType: "scene",
											},
										},
									},
								},
							},
							OnNoMotion: behavior_instance.OnNoMotion{
								After: behavior_instance.After{
									Minutes: 5,
								},
								RecallSingle: []behavior_instance.RecallSingleNoMotion{
									{
										Action: "all_off",
									},
								},
							},
							StartTime: behavior_instance.StartTime{
								Time: behavior_instance.Time{
									Hour:   23,
									Minute: 0,
								},
								Type: "time",
							},
						},
					},
				},
				Source: common.Reference{
					RID:   "f1fcd1cb-b439-4668-b6a9-0b57127af314",
					RType: "device",
				},
			},
			Dependees: []behavior_instance.Dependees{
				{
					Target: common.Reference{
						RID:   "f1fcd1cb-b439-4668-b6a9-0b57127af314",
						RType: "device",
					},
					Level: "critical",
					Type:  "ResourceDependee",
				},
				{
					Target: common.Reference{
						RID:   "0d960eab-68c6-4ed7-8c0d-a24ca756d58e",
						RType: "room",
					},
					Level: "critical",
					Type:  "ResourceDependee",
				},
				{
					Target: common.Reference{
						RID:   "68b39f81-1c15-4c82-bd0b-ab28606f3d2e",
						RType: "scene",
					},
					Level: "critical",
					Type:  "ResourceDependee",
				},
				{
					Target: common.Reference{
						RID:   "961975da-d821-4165-9bba-6f5495a9feef",
						RType: "temperature",
					},
					Level: "critical",
					Type:  "ResourceDependee",
				},
				{
					Target: common.Reference{
						RID:   "deb7e1cd-2bb8-46a8-a5c2-0f2ec387c235",
						RType: "light_level",
					},
					Level: "critical",
					Type:  "ResourceDependee",
				},
				{
					Target: common.Reference{
						RID:   "0ee6be79-9078-4352-ae1b-ce6d31f4cb8c",
						RType: "motion",
					},
					Level: "critical",
					Type:  "ResourceDependee",
				},
			},
		}
		var result common.ResourceList[behavior_instance.Data]
		err := json.Unmarshal(([]byte)(response), &result)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(expected, result.Data[0]); diff != "" {
			t.Errorf("Mismatch (-want +got):\n%s", diff)
		}
	})
}
