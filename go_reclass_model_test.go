package goreclass

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ModelTestSuite struct{ suite.Suite }

func (s *ModelTestSuite) TestDefault()                { s.compareFromYaml("default") }
func (s *ModelTestSuite) TestFeaturedSimple()         { s.compareFromYaml("featured/simple") }
func (s *ModelTestSuite) TestFeaturedWithClasses()    { s.compareFromYaml("featured/with-classes") }
func (s *ModelTestSuite) TestFeaturedWithNestedRefs() { s.compareFromYaml("featured/with-nested-refs") }
func (s *ModelTestSuite) TestFeaturedWithOverrides()  { s.compareFromYaml("featured/with-overrides") }
func (s *ModelTestSuite) TestFeaturedWithRefs()       { s.compareFromYaml("featured/with-refs") }
func (s *ModelTestSuite) TestFeaturedWithRelName()    { s.compareFromYaml("featured/with-relative-name") }
// FIXME: must work with class loop
//func (s *ModelTestSuite) TestExtendedWithClassLoop()  { s.compareFromYaml("extended/with-class-loop") }
func (s *ModelTestSuite) TestExtendedWithRefInClass() { s.compareFromYaml("extended/with-ref-in-class") }

func (s *ModelTestSuite) TestErroredOverrideConstant() {
	reclassFirstNodePath := path.Join("test/model/errored/override-constant/classes/first.yml")
	_, err := BuildInventory(reclassFirstNodePath)

	s.Error(err, "BuildInventory must fail if a constant is overrode.")
}

func TestModelSuiteTest(t *testing.T) { suite.Run(t, new(ModelTestSuite)) }
func (s *ModelTestSuite) compareFromYaml(modelName string) {
	reclassYamlPath := path.Join("test/model", modelName+".reclassed.yml")
	reclassYaml, err := ioutil.ReadFile(reclassYamlPath)
	if err != nil {
		s.FailNow(fmt.Sprintf("failed to read `reclassed` YAML file %s: %s", reclassYamlPath, err))
	}

	var reclassInventory Inventory
	if err := yaml.Unmarshal(reclassYaml, &reclassInventory); err != nil {
		s.FailNow(fmt.Sprintf("failed to unmarshal `reclassed` YAML file %s: %s", reclassYamlPath, err))
	}

	reclassFirstNodePath := path.Join("test/model", modelName, "classes/first.yml")
	goreclassInventory, err := BuildInventory(reclassFirstNodePath)
	if err != nil {
		s.FailNow(fmt.Sprintf("failed to build inventory from %s: %s", path.Join("test/model", modelName), err))
	}

	s.Equalf(&reclassInventory, goreclassInventory, "Inventories differs for model `%s`.", modelName)
}
