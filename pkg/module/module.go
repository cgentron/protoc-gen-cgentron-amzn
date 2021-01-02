package module

import (
	"bytes"
	"io"

	"github.com/cgentron/protoc-gen-cgentron-amzn/pkg/templates"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func New() pgs.Module { return &Module{ModuleBase: &pgs.ModuleBase{}} }

func (m *Module) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Module) Name() string { return "amzn" }

func (m *Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	buf := &bytes.Buffer{}

	for _, f := range targets {
		m.genAmzn(f, buf, m.Parameters())
	}

	return m.Artifacts()
}

func (m *Module) genAmzn(f pgs.File, buf *bytes.Buffer, params pgs.Parameters) {
	m.Push(f.Name().String())
	defer m.Pop()

	buf.Reset()
	v := initAmznVisitor(m, buf, "", params)
	m.CheckErr(pgs.Walk(v, f), "unable to generate AST tree")

	out := buf.String()

	m.AddGeneratorFile(
		m.ctx.OutputPath(f).SetExt(".amzn.go").String(),
		out,
	)
}

type AmznVisitor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	pgs.Parameters
	prefix string
	w      io.Writer
}

func initAmznVisitor(d pgs.DebuggerCommon, w io.Writer, prefix string, params pgs.Parameters) pgs.Visitor {
	p := AmznVisitor{
		prefix:         prefix,
		w:              w,
		Parameters:     params,
		DebuggerCommon: d,
	}

	p.Visitor = pgs.PassThroughVisitor(&p)

	return p
}

// VisitFile ...
func (p AmznVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
	tpl, err := templates.File(p.Parameters)
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(p.w, f)
	if err != nil {
		return nil, err
	}

	return p, err
}

// VisitService ...
func (p AmznVisitor) VisitService(s pgs.Service) (pgs.Visitor, error) {
	tpl, err := templates.Service(p.Parameters)
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(p.w, s)
	if err != nil {
		return nil, err
	}

	return p, err
}

// VisitMethod ...
func (p AmznVisitor) VisitMethod(m pgs.Method) (pgs.Visitor, error) {
	if m.ServerStreaming() {
		return p.visitMethodServerSideStreaming(m)
	}

	if m.ClientStreaming() {
		return p.visitMethodClientSideStreaming(m)
	}

	return p.visitMethod(m)
}

func (p AmznVisitor) visitMethod(m pgs.Method) (pgs.Visitor, error) {
	tpl, err := templates.Method(p.Parameters)
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(p.w, m)
	if err != nil {
		return nil, err
	}

	return p, err
}

func (p AmznVisitor) visitMethodServerSideStreaming(m pgs.Method) (pgs.Visitor, error) {
	tpl, err := templates.MethodServerStreaming(p.Parameters)
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(p.w, m)
	if err != nil {
		return nil, err
	}

	return p, err
}

func (p AmznVisitor) visitMethodClientSideStreaming(m pgs.Method) (pgs.Visitor, error) {
	tpl, err := templates.MethodClientStreaming(p.Parameters)
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(p.w, m)
	if err != nil {
		return nil, err
	}

	return p, err
}

// VisitMessage ...
func (p AmznVisitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
	tpl, err := templates.Message(p.Parameters)
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(p.w, m)
	if err != nil {
		return nil, err
	}

	return p, err
}

var _ pgs.Module = (*Module)(nil)
