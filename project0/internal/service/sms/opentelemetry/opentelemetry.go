package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"project0/internal/service/sms"
)

/*
Created by society-programmer on 2024/2/19.
*/


type Decorator struct {
	svc sms.Service
	tracer  trace.Tracer
}
// 关键是span的应用
func (d *Decorator) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//开启tracer,这里可以直接用tracer是因为，有初始化和设置全局traceProvider
	ctx, span := d.tracer.Start(ctx, "sms")
	defer span.End()
	span.SetAttributes(attribute.String("tpl",tplId))
	span.AddEvent("发短信")
	err := d.svc.Send(ctx,tplId,args)
	if err != nil {
		span.RecordError(err)
	}
	return err

}



