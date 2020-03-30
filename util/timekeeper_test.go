package util

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timeKeeper := NewTimeKeeper(map[string]string{"group": "main"})

	subKeeper := timeKeeper.Child(map[string]string{"subgroup": "C"})

	timer := timeKeeper.Start()
	timer2 := timeKeeper.Start()
	timer3 := subKeeper.Start()
	timer.Save(nil)
	time.Sleep(time.Millisecond * 200)
	timer2.Save(map[string]string{"postsleep": "1"})
	timer3.Save(map[string]string{"postsleep": "1", "wasnice": "yes"})

	records := timeKeeper.Records()
	recordsMatch := subKeeper.Records()
	require.Equal(t, records, recordsMatch)

	require.Equal(t, records[0].Tags, map[string]string{"group": "main"})
	require.True(t, records[0].Duration > time.Second*0)
	require.True(t, records[0].Duration < time.Millisecond*10)

	require.Equal(t, records[1].Tags, map[string]string{"group": "main", "postsleep": "1"})
	require.True(t, records[1].Duration > time.Millisecond*200)
	require.True(t, records[1].Duration < time.Millisecond*240)

	require.Equal(t, records[2].Tags, map[string]string{"group": "main", "postsleep": "1", "subgroup": "C", "wasnice": "yes"})
	require.True(t, records[2].Duration > time.Millisecond*200)
	require.True(t, records[2].Duration < time.Millisecond*240)

	require.True(t, len(records[0].String()) > 0)
}
