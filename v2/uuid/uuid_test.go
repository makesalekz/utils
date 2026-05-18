package uuid_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/makesalekz/utils/v2/uuid"
)

func TestNewFromActorIDPositive(t *testing.T) {
	var actorId int64 = 1378574
	uid := uuid.NewFromActorID(actorId)

	hexRepr := fmt.Sprintf("%x", actorId)

	reconstructedActorId := int64(
		binary.BigEndian.Uint64(
			[]byte{
				uid[7], uid[9], uid[10], uid[11], uid[12], uid[13], uid[14], uid[15],
			},
		),
	)

	t.Logf("%v, hex:%s", actorId, hexRepr)
	t.Logf("%v, hex:%x", reconstructedActorId, reconstructedActorId)
	t.Logf("uuid: %v", uid)

	require.Equal(t, fmt.Sprintf("%x", actorId), fmt.Sprintf("%x", reconstructedActorId))
}

func TestNewFromActorIDNegative(t *testing.T) {
	var actorId int64 = -1378574
	uid := uuid.NewFromActorID(actorId)

	hexRepr := "ffffffffffeaf6f2" // -1378574 in hex

	reconstructedActorId := int64(
		binary.BigEndian.Uint64(
			[]byte{
				uid[7], uid[9], uid[10], uid[11], uid[12], uid[13], uid[14], uid[15],
			},
		),
	)

	t.Logf("%v, hex:%s", actorId, hexRepr)
	t.Logf("%v, hex:%x", reconstructedActorId, reconstructedActorId)
	t.Logf("uuid: %v", uid)

	require.Equal(t, fmt.Sprintf("%x", actorId), fmt.Sprintf("%x", reconstructedActorId))
}

func TestNewFromActorIDZero(t *testing.T) {
	var actorId int64 = 0
	uid := uuid.NewFromActorID(actorId)

	hexRepr := "0000000000000000" // 0 in hex

	reconstructedActorId := int64(
		binary.BigEndian.Uint64(
			[]byte{
				uid[7], uid[9], uid[10], uid[11], uid[12], uid[13], uid[14], uid[15],
			},
		),
	)

	t.Logf("%v, hex:%s", actorId, hexRepr)
	t.Logf("%v, hex:%x", reconstructedActorId, reconstructedActorId)
	t.Logf("uuid: %v", uid)

	require.Equal(t, fmt.Sprintf("%x", actorId), fmt.Sprintf("%x", reconstructedActorId))
}
