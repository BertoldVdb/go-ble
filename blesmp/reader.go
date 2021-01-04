package blesmp

import "context"

func (c *SMPConn) reader() {
	ctx := context.Background()

	defer c.conn.UseDone()
	defer close(c.pduRx)
	defer c.conn.Close()
	c.conn.UseStart()

	for {
		pdu, err := c.conn.ReadBuffer(ctx)
		if err != nil {
			c.logger.WithError(err).Debug("PSM Read error")
			return
		}

		c.pduRx <- pdu
	}
}
