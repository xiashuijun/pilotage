package session

import (
	"sync"

	"github.com/abiosoft/ishell"
	"github.com/HardySimpson/pilotage/pkg/vfs"
	"k8s.io/client-go/kubernetes"
	"strings"
)


type Session struct {
	*ishell.Shell
	kubecli kubernetes.Interface
	path string
	vfs *vfs.VFS
	currentNode *vfs.Node
	lock sync.Mutex
}

func (s *Session) cd(c *ishell.Context) {
	args := c.Args

	if len(args) > 1 {
		c.Println("pilotage cd: too many arguments")
	} else if len(args) == 0 {
		s.currentNode = s.vfs.RootNode
	} else if args[0] == ".." {
		if s.currentNode.Parent != nil {
			s.currentNode = s.currentNode.Parent
		}
	} else if args[0] == "." {
		//do nothing
		return
	}

	n := s.currentNode.GetChild(args[0])
	if n != nil {
		s.currentNode = n
	}

	s.SetPrompt(s.currentNode.Path + "$ ")

}

func (s *Session) inspect(c *ishell.Context) {
	s.Println(s.currentNode)
}

func (s *Session) cdCompleter(args []string) []string {
	return s.currentNode.ListChildrenName()

}

func (s *Session) pwd(c *ishell.Context) {
	s.Shell.Println(s.currentNode.Path)
}

func (s *Session) ls(c *ishell.Context) {
	/*
	if s.currentNode.Fresh != nil {
		s.currentNode.Fresh()
	}
	*/
	s.Shell.Println(strings.Join(s.currentNode.ListChildrenName(), " "))
}


func New(k kubernetes.Interface) *Session {
	s := &Session{
		Shell: ishell.New(),
		kubecli: k,
		vfs: vfs.NewVFS(),
	}
	s.currentNode = s.vfs.RootNode

	s.AddCmd(&ishell.Cmd{
		Name: "cd",
		Help: "change to a path",
		Func: s.cd,
		Completer: s.cdCompleter,
	})
	s.AddCmd(&ishell.Cmd{
		Name: "ls",
		Help: "list sub object",
		Func: s.ls,
	})
	s.AddCmd(&ishell.Cmd{
		Name: "pwd",
		Help: "show path",
		Func: s.pwd,
	})
	s.AddCmd(&ishell.Cmd{
		Name: "inspect",
		Help: "inspect current vfs.Node",
		Func: s.inspect,
	})

	s.SetPrompt("/$ ")

	addDefaultTree(s.vfs.RootNode, s.kubecli)


	return s
}
