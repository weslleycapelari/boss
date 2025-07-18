package compiler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/weslleycapelari/boss/pkg/compiler/graphs"
	"github.com/weslleycapelari/boss/pkg/consts"
	"github.com/weslleycapelari/boss/pkg/env"
	"github.com/weslleycapelari/boss/pkg/models"
	"github.com/weslleycapelari/boss/pkg/msg"
	"github.com/weslleycapelari/boss/utils"
)

func Build(pkg *models.Package) {
	buildOrderedPackages(pkg)
	graph := LoadOrderGraphAll(pkg)
	saveLoadOrder(graph)
}

func saveLoadOrder(queue *graphs.NodeQueue) {
	var projects = ""
	for {
		if queue.IsEmpty() {
			break
		}
		node := queue.Dequeue()
		dependencyPath := filepath.Join(env.GetModulesDir(), node.Dep.Name(), consts.FilePackage)
		if dependencyPackage, err := models.LoadPackageOther(dependencyPath); err == nil {
			for _, value := range dependencyPackage.Projects {
				projects += strings.TrimSuffix(filepath.Base(value), filepath.Ext(value)) + consts.FileExtensionBpl + "\n"
			}
		}
	}
	outDir := filepath.Join(env.GetModulesDir(), consts.BplFolder, consts.FileBplOrder)

	utils.HandleError(os.WriteFile(outDir, []byte(projects), 0600))
}

func buildOrderedPackages(pkg *models.Package) {
	pkg.Lock.Save()
	queue := loadOrderGraph(pkg)
	for {
		if queue.IsEmpty() {
			break
		}
		node := queue.Dequeue()
		dependencyPath := filepath.Join(env.GetModulesDir(), node.Dep.Name())

		dependency := pkg.Lock.GetInstalled(node.Dep)

		msg.Info("Building %s", node.Dep.Name())
		dependency.Changed = false
		if dependencyPackage, err := models.LoadPackageOther(filepath.Join(dependencyPath, consts.FilePackage)); err == nil {
			dprojs := dependencyPackage.Projects
			if len(dprojs) > 0 {
				for _, dproj := range dprojs {
					dprojPath, _ := filepath.Abs(filepath.Join(env.GetModulesDir(), node.Dep.Name(), dproj))
					if !compile(dprojPath, &node.Dep, pkg.Lock) {
						dependency.Failed = true
					}
				}
				ensureArtifacts(&dependency, node.Dep, env.GetModulesDir())
				moveArtifacts(node.Dep, env.GetModulesDir())
			}
		}
		pkg.Lock.SetInstalled(node.Dep, dependency)
	}
}
