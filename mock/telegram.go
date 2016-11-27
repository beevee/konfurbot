// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/beevee/konfurbot/telegram (interfaces: TelebotInterface)

package mock

import (
	time "time"

	gomock "github.com/golang/mock/gomock"
	telebot "github.com/tucnak/telebot"
)

// Mock of TelebotInterface interface
type MockTelebotInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockTelebotInterfaceRecorder
}

// Recorder for MockTelebotInterface (not exported)
type _MockTelebotInterfaceRecorder struct {
	mock *MockTelebotInterface
}

func NewMockTelebotInterface(ctrl *gomock.Controller) *MockTelebotInterface {
	mock := &MockTelebotInterface{ctrl: ctrl}
	mock.recorder = &_MockTelebotInterfaceRecorder{mock}
	return mock
}

func (_m *MockTelebotInterface) EXPECT() *_MockTelebotInterfaceRecorder {
	return _m.recorder
}

func (_m *MockTelebotInterface) Listen(_param0 chan telebot.Message, _param1 time.Duration) {
	_m.ctrl.Call(_m, "Listen", _param0, _param1)
}

func (_mr *_MockTelebotInterfaceRecorder) Listen(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Listen", arg0, arg1)
}

func (_m *MockTelebotInterface) SendMessage(_param0 telebot.Recipient, _param1 string, _param2 *telebot.SendOptions) error {
	ret := _m.ctrl.Call(_m, "SendMessage", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockTelebotInterfaceRecorder) SendMessage(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SendMessage", arg0, arg1, arg2)
}
