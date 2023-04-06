package grpcutil

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func StreamToReader[T proto.Message](ctx context.Context, stream grpc.ServerStream) io.Reader {
	r, w := io.Pipe()

	var processMessage = func() error {
		if err := ctx.Err(); err != nil {
			return err
		}

		var message T

		if err := stream.RecvMsg(&message); err != nil {
			return err
		}

		data, err := proto.Marshal(message)
		if err != nil {
			return err
		}

		_, err = w.Write(data)

		return err
	}

	go func() {
		for {
			if err := processMessage(); err != nil {
				if !errors.Is(err, io.EOF) {
					w.CloseWithError(err)
				}
				w.Close()
				return
			}
		}
	}()

	return r
}

func StreamToWriter[T proto.Message](ctx context.Context, stream grpc.ServerStream) io.Writer {
	r, w := io.Pipe()

	var processMessage = func() error {
		if err := ctx.Err(); err != nil {
			return err
		}

		buffer := make([]byte, 1024)

		n, err := r.Read(buffer)
		if err != nil {
			return err
		}

		var message T

		if err = proto.Unmarshal(buffer[:n], message); err != nil {
			return err
		}

		return stream.SendMsg(&message)
	}

	go func() {
		for {
			if err := processMessage(); err != nil {
				if !errors.Is(err, io.EOF) {
					w.CloseWithError(err)
				}
				w.Close()
				return
			}
		}
	}()

	return w
}
