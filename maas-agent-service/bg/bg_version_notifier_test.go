package bg

import (
	"context"
	"testing"
	"time"

	"github.com/netcracker/qubership-core-maas-agent/maas-agent-service/v2/model"

	"github.com/stretchr/testify/assert"
)

func TestIsContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	assert.False(t, IsContextCancelled(ctx))

	cancel()
	assert.True(t, IsContextCancelled(ctx))
}

func Test_cancelableSleep(t *testing.T) {
	start := time.Now()
	CancelableSleep(context.Background(), 1*time.Second)

	if time.Now().Sub(start) < 1*time.Second {
		assert.Fail(t, "sleep must be cancelled")
	}
}

func Test_cancelableSleepCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	CancelableSleep(ctx, 1*time.Second)

	if time.Now().Sub(start) >= 1*time.Second {
		assert.Fail(t, "sleep must be cancelled")
	}
}

func TestApplyVersionsChangeInit(t *testing.T) {
	event := model.CpWatcherMessageDto{State: model.CpVersionsDto{
		{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
	}}
	assert.Equal(t, event.State, applyVersionChange(context.TODO(), nil, event))
}

func TestApplyVersionsChangeAddCandidate(t *testing.T) {
	versions := model.CpVersionsDto{
		{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
	}

	event := model.CpWatcherMessageDto{State: nil,
		Changes: []model.CpChange{
			{New: &model.CpDeploymentVersion{Version: "v2", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T19:32:45.501289713Z", UpdatedWhen: "2022-11-09T19:32:45.501289856Z"}},
		},
	}

	expected := model.CpVersionsDto{
		{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
		{Version: "v2", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T19:32:45.501289713Z", UpdatedWhen: "2022-11-09T19:32:45.501289856Z"},
	}
	assert.Equal(t, expected, applyVersionChange(context.TODO(), versions, event))
}
func TestApplyVersionsChangePromote(t *testing.T) {
	versions := model.CpVersionsDto{
		{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
		{Version: "v2", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T19:32:45.501289713Z", UpdatedWhen: "2022-11-09T19:32:45.501289856Z"},
	}

	event := model.CpWatcherMessageDto{State: nil,
		Changes: []model.CpChange{
			{
				Old: &model.CpDeploymentVersion{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
				New: &model.CpDeploymentVersion{Version: "v1", Stage: "LEGACY", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
			},
			{
				Old: &model.CpDeploymentVersion{Version: "v2", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T19:32:45.501289713Z", UpdatedWhen: "2022-11-09T19:32:45.501289856Z"},
				New: &model.CpDeploymentVersion{Version: "v2", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:32:45.501289713Z", UpdatedWhen: "2022-11-09T19:32:45.501289856Z"},
			},
		},
	}

	expected := model.CpVersionsDto{
		{Version: "v1", Stage: "LEGACY", CreatedWhen: "2022-11-09T19:17:13.000165Z", UpdatedWhen: "2022-11-09T19:17:13.000165Z"},
		{Version: "v2", Stage: "ACTIVE", CreatedWhen: "2022-11-09T19:32:45.501289713Z", UpdatedWhen: "2022-11-09T19:32:45.501289856Z"},
	}
	assert.Equal(t, expected, applyVersionChange(context.TODO(), versions, event))
}

func TestApplyVersionsChangeRollback(t *testing.T) {
	versions := model.CpVersionsDto{
		{Version: "v1", Stage: "LEGACY", CreatedWhen: "2022-11-09T20:25:46.808564Z", UpdatedWhen: "2022-11-09T20:25:46.808564Z"},
		{Version: "v2", Stage: "ACTIVE", CreatedWhen: "2022-11-09T20:38:46.167563449Z", UpdatedWhen: "2022-11-09T20:38:46.16756354Z"},
		{Version: "v3", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T20:44:57.374169286Z", UpdatedWhen: "2022-11-09T20:44:57.374169369Z"},
	}

	event := model.CpWatcherMessageDto{State: nil,
		Changes: []model.CpChange{
			{
				Old: &model.CpDeploymentVersion{Version: "v3", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T20:44:57.374169286Z", UpdatedWhen: "2022-11-09T20:44:57.374169369Z"},
				New: nil,
			},
			{
				Old: &model.CpDeploymentVersion{Version: "v1", Stage: "LEGACY", CreatedWhen: "2022-11-09T20:25:46.808564Z", UpdatedWhen: "2022-11-09T20:25:46.808564Z"},
				New: &model.CpDeploymentVersion{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T20:25:46.808564Z", UpdatedWhen: "2022-11-09T20:25:46.808564Z"},
			},
			{
				Old: &model.CpDeploymentVersion{Version: "v2", Stage: "ACTIVE", CreatedWhen: "2022-11-09T20:38:46.167563449Z", UpdatedWhen: "2022-11-09T20:38:46.16756354Z"},
				New: &model.CpDeploymentVersion{Version: "v2", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T20:38:46.167563449Z", UpdatedWhen: "2022-11-09T20:38:46.16756354Z"},
			},
		},
	}

	expected := model.CpVersionsDto{
		{Version: "v1", Stage: "ACTIVE", CreatedWhen: "2022-11-09T20:25:46.808564Z", UpdatedWhen: "2022-11-09T20:25:46.808564Z"},
		{Version: "v2", Stage: "CANDIDATE", CreatedWhen: "2022-11-09T20:38:46.167563449Z", UpdatedWhen: "2022-11-09T20:38:46.16756354Z"},
	}

	assert.Equal(t, expected, applyVersionChange(context.TODO(), versions, event))
}
