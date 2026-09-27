package huhu1

/*
================================================================================
문제 4: 자율주행 차량 상태 머신
================================================================================

난이도: Hard
주제: 상태 머신, 동시성, 이벤트 기반 아키텍처
회사: 42dot (자율주행/모빌리티)

【 문제 설명 】
자율주행 차량의 상태를 관리하는 유한 상태 머신(FSM)을 구현하세요.
차량은 다양한 상태 간 전이를 하며, 각 전이는 조건과 액션을 가집니다.

상태 (State):
- PARKED: 주차됨
- IDLE: 대기 (시동 ON, 정차)
- CRUISING: 순항 중
- LANE_CHANGING: 차선 변경 중
- TURNING: 회전 중
- EMERGENCY_STOP: 비상 정지
- CHARGING: 충전 중

이벤트 (Event):
- START: 시동 시작
- STOP: 시동 정지
- ACCELERATE: 가속
- BRAKE: 감속/정지
- LANE_CHANGE_START: 차선 변경 시작
- LANE_CHANGE_COMPLETE: 차선 변경 완료
- TURN_START: 회전 시작
- TURN_COMPLETE: 회전 완료
- EMERGENCY: 비상 상황
- RESUME: 복귀
- CHARGE_START: 충전 시작
- CHARGE_COMPLETE: 충전 완료

요구사항:
1. Transition(event): 상태 전이 실행
2. GetCurrentState(): 현재 상태 반환
3. OnEnter(state, handler): 상태 진입 시 콜백
4. OnExit(state, handler): 상태 퇴장 시 콜백
5. CanTransition(event): 전이 가능 여부 확인
6. GetStateHistory(): 상태 변경 이력 반환

【 성능 요구사항 】
- 전이 연산: O(1)
- 동시 접근 안전
- 상태 히스토리: 최근 100개 유지

【 인터페이스 】
type VehicleState string
type VehicleEvent string

type StateTransition struct {
    From   VehicleState
    Event  VehicleEvent
    To     VehicleState
    Guard  func() bool       // 전이 조건 (nil이면 항상 허용)
    Action func()            // 전이 시 실행할 액션
}

type StateHistory struct {
    From      VehicleState
    To        VehicleState
    Event     VehicleEvent
    Timestamp time.Time
}

type VehicleFSM interface {
    GetCurrentState() VehicleState
    Transition(event VehicleEvent) error
    CanTransition(event VehicleEvent) bool
    OnEnter(state VehicleState, handler func())
    OnExit(state VehicleState, handler func())
    GetStateHistory() []StateHistory
}

【 힌트 】
1. 전이 테이블: map[(state, event)] -> nextState
2. sync.RWMutex로 동시성 제어
3. Ring buffer로 히스토리 관리
4. Observer 패턴으로 콜백 관리

【 실무 연관성 】
- 자율주행 차량 모드 관리
- 로봇 제어 시스템
- 게임 AI 상태 관리
================================================================================
*/

import "time"

// VehicleState 차량 상태
type VehicleState string

const (
	PARKED         VehicleState = "PARKED"
	IDLE           VehicleState = "IDLE"
	CRUISING       VehicleState = "CRUISING"
	LANE_CHANGING  VehicleState = "LANE_CHANGING"
	TURNING        VehicleState = "TURNING"
	EMERGENCY_STOP VehicleState = "EMERGENCY_STOP"
	CHARGING       VehicleState = "CHARGING"
)

// VehicleEvent 차량 이벤트
type VehicleEvent string

const (
	START                VehicleEvent = "START"
	STOP                 VehicleEvent = "STOP"
	ACCELERATE           VehicleEvent = "ACCELERATE"
	BRAKE                VehicleEvent = "BRAKE"
	LANE_CHANGE_START    VehicleEvent = "LANE_CHANGE_START"
	LANE_CHANGE_COMPLETE VehicleEvent = "LANE_CHANGE_COMPLETE"
	TURN_START           VehicleEvent = "TURN_START"
	TURN_COMPLETE        VehicleEvent = "TURN_COMPLETE"
	EVENT_EMERGENCY      VehicleEvent = "EMERGENCY"
	RESUME               VehicleEvent = "RESUME"
	CHARGE_START         VehicleEvent = "CHARGE_START"
	CHARGE_COMPLETE      VehicleEvent = "CHARGE_COMPLETE"
)

// StateTransition 상태 전이 정의
type StateTransition struct {
	From   VehicleState
	Event  VehicleEvent
	To     VehicleState
	Guard  func() bool // 전이 조건
	Action func()      // 전이 액션
}

// StateHistory 상태 변경 이력
type StateHistory struct {
	From      VehicleState
	To        VehicleState
	Event     VehicleEvent
	Timestamp time.Time
}

// VehicleFSM 차량 상태 머신 인터페이스
type VehicleFSM interface {
	GetCurrentState() VehicleState
	Transition(event VehicleEvent) error
	CanTransition(event VehicleEvent) bool
	OnEnter(state VehicleState, handler func())
	OnExit(state VehicleState, handler func())
	GetStateHistory() []StateHistory
}

// VehicleFSMImpl 구현체
type VehicleFSMImpl struct {
	// 여기에 필드를 정의하세요
}

// NewVehicleFSM 생성자
func NewVehicleFSM(initial VehicleState) *VehicleFSMImpl {
	// 여기에 코드를 작성하세요
	return nil
}

// RegisterTransition 전이 규칙 등록
func (f *VehicleFSMImpl) RegisterTransition(t StateTransition) {
	// 여기에 코드를 작성하세요
}

// GetCurrentState 현재 상태 반환
func (f *VehicleFSMImpl) GetCurrentState() VehicleState {
	// 여기에 코드를 작성하세요
	return ""
}

// Transition 상태 전이 실행
func (f *VehicleFSMImpl) Transition(event VehicleEvent) error {
	// 여기에 코드를 작성하세요
	return nil
}

// CanTransition 전이 가능 여부 확인
func (f *VehicleFSMImpl) CanTransition(event VehicleEvent) bool {
	// 여기에 코드를 작성하세요
	return false
}

// OnEnter 상태 진입 콜백 등록
func (f *VehicleFSMImpl) OnEnter(state VehicleState, handler func()) {
	// 여기에 코드를 작성하세요
}

// OnExit 상태 퇴장 콜백 등록
func (f *VehicleFSMImpl) OnExit(state VehicleState, handler func()) {
	// 여기에 코드를 작성하세요
}

// GetStateHistory 상태 변경 이력 반환
func (f *VehicleFSMImpl) GetStateHistory() []StateHistory {
	// 여기에 코드를 작성하세요
	return nil
}
