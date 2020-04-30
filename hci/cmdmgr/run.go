package hcicmdmgr

type HCICommand struct {
	OGF    int
	OCF    int
	Params []byte
}

type QueueIndex int

func (s *CommandManager) CommandRunSync(queue QueueIndex, cmd HCICommand, output []byte) ([]byte, error) {
	return s.queues[queue].commandRun(cmd, output, true, nil)
}

func (s *CommandManager) CommandRunAsync(queue QueueIndex, cmd HCICommand, cb CommandCompleteCallback) error {
	_, err := s.queues[queue].commandRun(cmd, nil, false, cb)
	return err
}

func (s *CommandManager) CommandRun(queue QueueIndex, cmd HCICommand) error {
	_, err := s.queues[queue].commandRun(cmd, nil, false, nil)
	return err
}

type HCICommandBuffer struct {
	token  *commandToken
	queue  QueueIndex
	Buffer []byte
}

func (s *CommandManager) CommandRunGetBuffer(queue QueueIndex, cmd HCICommand, cb CommandCompleteCallback) (HCICommandBuffer, error) {
	token, err := s.queues[queue].commandRunGetToken(cmd, cb == nil, cb)
	return HCICommandBuffer{
		token:  token,
		queue:  queue,
		Buffer: token.data,
	}, err
}

func (s *CommandManager) CommandRunPutBuffer(buffer HCICommandBuffer) ([]byte, error) {
	buffer.token.data = buffer.Buffer
	return s.queues[buffer.queue].commandRunPutToken(buffer.token)
}

func (s *CommandManager) CommandRunReleaseBuffer(buffer HCICommandBuffer) error {
	return s.queues[buffer.queue].commandRunReleaseToken(buffer.token)
}
