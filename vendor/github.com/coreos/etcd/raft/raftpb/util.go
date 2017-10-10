package raftpb

import (
	"fmt"
	"strings"
	"encoding/gob"
	"bytes"
	"log"
)

func (m *Message) BeautiString() string {
	return fmt.Sprintf("{%s %d<-%d term:%d idx:%d log%d, ents:%s cmt:%d snap:%s reject:%v hint:%d ctx:%s}", m.Type.String()[3:], m.To, m.From, m.Term, m.Index, m.LogTerm, EntriesStr(m.Entries), m.Commit, m.Snapshot.BeautiString(), m.Reject, m.RejectHint, string(m.Context))
}

func (cc *ConfChange) BeautiString() string {
	return fmt.Sprintf("{%v %d %d %s}", cc.Type.String()[10:], cc.ID, cc.NodeID, string(cc.Context))
}

func (cs *ConfState) BeautiString() string {
	return fmt.Sprintf("%v", cs.Nodes)
}

func (sm *SnapshotMetadata) BeautiString() string {
	return fmt.Sprintf("{confstat:%s idx:%d term:%d}", sm.ConfState.BeautiString(), sm.Index, sm.Term)
}

func (s *Snapshot) BeautiString() string {
	return fmt.Sprintf("{meta:%s data:%s}", s.Metadata.BeautiString(), string(s.Data))
}

type kv struct {
	Key string
	Val string
}

func EntriesStr(entries []Entry) string {
	var arrs []string
	for i := range entries {
		if entries[i].Type == EntryConfChange {
			var conf ConfChange
			if err := conf.Unmarshal(entries[i].Data); err != nil {
				log.Fatalf("raftexample: could not decode conf (%v)", err)
			}
			arrs = append(arrs, fmt.Sprintf("c term:%d idx:%d data:%s", entries[i].Term, entries[i].Index, conf.BeautiString()))
		} else {
			var kv kv
			if err := gob.NewDecoder(bytes.NewBuffer(entries[i].Data)).Decode(&kv); err != nil {
				arrs = append(arrs, fmt.Sprintf("d term:%d idx:%d data:%s", entries[i].Term, entries[i].Index, string(entries[i].Data)))
			} else {
				arrs = append(arrs, fmt.Sprintf("d term:%d idx:%d data:%s=%s", entries[i].Term, entries[i].Index, kv.Key, kv.Val))
			}
		}
	}
	return "{"+strings.Join(arrs, ",")+"}"
}

func MessagesStr(msgs []Message) string {
	var arr []string
	for i := range msgs {
		arr = append(arr, msgs[i].BeautiString())
	}
	return "{"+strings.Join(arr, ",")+"}"
}
