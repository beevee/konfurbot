// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/beevee/konfurbot (interfaces: ScheduleStorage)

package mock

import (
	konfurbot "github.com/beevee/konfurbot"
	gomock "github.com/golang/mock/gomock"
	time "time"
)

// Mock of ScheduleStorage interface
type MockScheduleStorage struct {
	ctrl     *gomock.Controller
	recorder *_MockScheduleStorageRecorder
}

// Recorder for MockScheduleStorage (not exported)
type _MockScheduleStorageRecorder struct {
	mock *MockScheduleStorage
}

func NewMockScheduleStorage(ctrl *gomock.Controller) *MockScheduleStorage {
	mock := &MockScheduleStorage{ctrl: ctrl}
	mock.recorder = &_MockScheduleStorageRecorder{mock}
	return mock
}

func (_m *MockScheduleStorage) EXPECT() *_MockScheduleStorageRecorder {
	return _m.recorder
}

func (_m *MockScheduleStorage) AddEvent(_param0 konfurbot.Event) {
	_m.ctrl.Call(_m, "AddEvent", _param0)
}

func (_mr *_MockScheduleStorageRecorder) AddEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddEvent", arg0)
}

func (_m *MockScheduleStorage) GetCurrentEventsByType(_param0 string, _param1 time.Time) []konfurbot.Event {
	ret := _m.ctrl.Call(_m, "GetCurrentEventsByType", _param0, _param1)
	ret0, _ := ret[0].([]konfurbot.Event)
	return ret0
}

func (_mr *_MockScheduleStorageRecorder) GetCurrentEventsByType(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCurrentEventsByType", arg0, arg1)
}

func (_m *MockScheduleStorage) GetDayEventsByType(_param0 string) []konfurbot.Event {
	ret := _m.ctrl.Call(_m, "GetDayEventsByType", _param0)
	ret0, _ := ret[0].([]konfurbot.Event)
	return ret0
}

func (_mr *_MockScheduleStorageRecorder) GetDayEventsByType(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDayEventsByType", arg0)
}

func (_m *MockScheduleStorage) GetEventsByType(_param0 string) []konfurbot.Event {
	ret := _m.ctrl.Call(_m, "GetEventsByType", _param0)
	ret0, _ := ret[0].([]konfurbot.Event)
	return ret0
}

func (_mr *_MockScheduleStorageRecorder) GetEventsByType(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEventsByType", arg0)
}

func (_m *MockScheduleStorage) GetEventsByTypeAndSubtype(_param0 string, _param1 string) []konfurbot.Event {
	ret := _m.ctrl.Call(_m, "GetEventsByTypeAndSubtype", _param0, _param1)
	ret0, _ := ret[0].([]konfurbot.Event)
	return ret0
}

func (_mr *_MockScheduleStorageRecorder) GetEventsByTypeAndSubtype(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEventsByTypeAndSubtype", arg0, arg1)
}

func (_m *MockScheduleStorage) GetNextEventsByType(_param0 string, _param1 time.Time, _param2 time.Duration) []konfurbot.Event {
	ret := _m.ctrl.Call(_m, "GetNextEventsByType", _param0, _param1, _param2)
	ret0, _ := ret[0].([]konfurbot.Event)
	return ret0
}

func (_mr *_MockScheduleStorageRecorder) GetNextEventsByType(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetNextEventsByType", arg0, arg1, arg2)
}

func (_m *MockScheduleStorage) GetNightEventsByType(_param0 string) []konfurbot.Event {
	ret := _m.ctrl.Call(_m, "GetNightEventsByType", _param0)
	ret0, _ := ret[0].([]konfurbot.Event)
	return ret0
}

func (_mr *_MockScheduleStorageRecorder) GetNightEventsByType(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetNightEventsByType", arg0)
}

func (_m *MockScheduleStorage) SetNightCutoff(_param0 time.Time) {
	_m.ctrl.Call(_m, "SetNightCutoff", _param0)
}

func (_mr *_MockScheduleStorageRecorder) SetNightCutoff(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetNightCutoff", arg0)
}
