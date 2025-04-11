package cryo

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

// ScheduledTaskType 定时任务类型别名
type ScheduledTaskType int

// ScheduledTaskStatus 定时任务状态别名
type ScheduledTaskStatus string

const (
	DelayTaskType    ScheduledTaskType = iota // 延迟任务，n秒后执行
	IntervalTaskType                          // 间隔任务，每n秒执行一次
	CronTaskType                              // 定时任务，可以指定具体时间来执行
)

const (
	TaskPending ScheduledTaskStatus = "pending" // 任务未开始
	TaskRunning ScheduledTaskStatus = "running" // 任务正在执行
	TaskFailed  ScheduledTaskStatus = "failed"  // 任务执行失败
	TaskStopped ScheduledTaskStatus = "stopped" // 任务已停止
)

// ScheduledTask 定时任务结构体
type ScheduledTask struct {
	id            string              // 任务ID
	name          string              // 任务名称
	taskType      ScheduledTaskType   // 任务类型
	duration      time.Duration       // 任务执行间隔，单位 秒
	cron          string              // 定时任务的cron表达式
	isWithSeconds bool                // 是否包含秒数
	isInstantly   bool                // 是否立即执行
	status        ScheduledTaskStatus // 任务状态
	err           error               // 如果任务执行时出现异常，则会在这里保存异常信息

	Task func() error // 任务对象，仅在使用gocron进行调度时使用
	Job  gocron.Job   // gocron任务对象
}

// GetTaskId 获取任务ID
func (st *ScheduledTask) GetTaskId() string {
	return st.id
}

// GetTaskName 获取任务名称
func (st *ScheduledTask) GetTaskName() string {
	return st.name
}

// GetTaskType 获取任务类型
func (st *ScheduledTask) GetTaskType() ScheduledTaskType {
	return st.taskType
}

// GetDuration 获取任务执行间隔
func (st *ScheduledTask) GetDuration() time.Duration {
	return st.duration
}

// GetCron 获取任务的cron表达式
func (st *ScheduledTask) GetCron() string {
	return st.cron
}

// GetStatus 获取任务状态
func (st *ScheduledTask) GetStatus() ScheduledTaskStatus {
	return st.status
}

// Set 向调度器设置任务
func (st *ScheduledTask) Set(b *Bot) {
	// 设置任务状态为正在执行
	st.status = TaskRunning

	if st.taskType == DelayTaskType {
		// 给任务嵌套一个包装器用于设置任务状态
		// 任务执行完成后设置任务状态为完成
		task := func() {
			err := st.Task()
			if err != nil {
				st.status = TaskFailed
				st.err = err
				// 发送任务失败事件
				SendScheduledTaskFailedEvent(b, st)
			}

			st.status = TaskStopped
			// 发送任务成功事件
			SendScheduledTaskSuccessEvent(b, st)
			// 发送任务停止事件
			SendScheduledTaskStoppedEvent(b, st)
		}

		job, _ := b.scheduler.NewJob(
			gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTime(time.Now().Add(time.Duration(st.duration))), // 设置触发时间
			),
			gocron.NewTask(task),
		)
		st.Job = job // 传递任务对象
	}

	if st.taskType == IntervalTaskType {
		task := func() {
			err := st.Task()
			if err != nil {
				st.status = TaskFailed
				st.err = err
				// 发送任务失败事件
				SendScheduledTaskFailedEvent(b, st)
			} else {
				// 发送任务成功事件
				SendScheduledTaskSuccessEvent(b, st)
			}
		}

		if st.isInstantly {
			job, _ := b.scheduler.NewJob(
				gocron.DurationJob(time.Duration(st.duration)), // 设置执行间隔[^1]
				gocron.NewTask(task),
				gocron.WithStartAt(gocron.WithStartImmediately()), // 设置立即执行
			)
			st.Job = job // 传递任务对象
		} else {
			job, _ := b.scheduler.NewJob(
				gocron.DurationJob(time.Duration(st.duration)), // 设置执行间隔[^1]
				gocron.NewTask(task),
			)
			st.Job = job // 传递任务对象
		}
	}

	if st.taskType == CronTaskType {
		task := func() {
			err := st.Task()
			if err != nil {
				st.status = TaskFailed
				st.err = err
				// 发送任务失败事件
				SendScheduledTaskFailedEvent(b, st)
			} else {
				// 发送任务成功事件
				SendScheduledTaskSuccessEvent(b, st)
			}
		}

		job, _ := b.scheduler.NewJob(
			gocron.CronJob(st.cron, st.isWithSeconds), // 设置执行时间
			gocron.NewTask(task),
		)
		st.Job = job // 传递任务对象
	}

	// 发送任务注册事件
	SendScheduledTaskRegisteredEvent(b, st)
}

// Stop 停止任务
func (st *ScheduledTask) Stop(b *Bot) error {
	if st.Job != nil {
		id := st.Job.ID()
		err := b.scheduler.RemoveJob(id)
		if err != nil {
			return err
		} // 停止任务
		st.status = TaskStopped
		st.Job = nil // 清空任务对象

		// 发送任务停止事件
		SendScheduledTaskStoppedEvent(b, st)

		return nil
	}
	return nil
}

// GetError 获取任务执行时的错误信息
func (st *ScheduledTask) GetError() error {
	return st.err
}

// NewDelayTask 创建一个新的延迟任务
func NewDelayTask(name string, duration time.Duration, task func() error) *ScheduledTask {
	return &ScheduledTask{
		id:            newUUID(),
		name:          name,
		taskType:      DelayTaskType,
		duration:      duration,
		cron:          "",
		isWithSeconds: false,
		isInstantly:   false,
		status:        TaskPending,
		Task:          task,
		Job:           nil,
	}
}

// NewIntervalTask 创建一个新的间隔任务
func NewIntervalTask(name string, duration time.Duration, isInstantly bool, task func() error) *ScheduledTask {
	return &ScheduledTask{
		id:            newUUID(),
		name:          name,
		taskType:      IntervalTaskType,
		duration:      duration,
		cron:          "",
		isWithSeconds: false,
		isInstantly:   isInstantly,
		status:        TaskPending,
		Task:          task,
		Job:           nil,
	}
}

// NewCronTask 创建一个新的定时任务
func NewCronTask(name string, cron string, isWithSeconds bool, task func() error) *ScheduledTask {
	return &ScheduledTask{
		id:            newUUID(),
		name:          name,
		taskType:      CronTaskType,
		duration:      0,
		cron:          cron,
		isWithSeconds: isWithSeconds,
		isInstantly:   false,
		status:        TaskPending,
		Task:          task,
		Job:           nil,
	}
}
