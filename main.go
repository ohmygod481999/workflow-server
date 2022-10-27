package main

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"callbot/workflow/auth"
	workflow_2 "callbot/workflow/workflow-2"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()

	viper.SetConfigName("settings")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	authCfg := auth.Config{
		PassportUri:        viper.GetString("passport_uri"),
		VintalkServicesUri: viper.GetString("vintalk_services"),
	}

	{
		app.Use(auth.NewPassportAuthenticator(authCfg))
		app.Use(auth.NewCallbotAgentAuthenticator(authCfg))
	}

	argoCfg := argo_adapter.Config{
		ArgoUri: viper.GetString("argo_uri"),
	}
	argoAdapter := argo_adapter.New(argoCfg)

	// auth.RequireUserOrCallbotAgent
	app.Get("/", func(c *fiber.Ctx) error {
		workflows := argoAdapter.GetWorkflows()
		return c.JSON(workflows)
	})

	app.Get("/wf-templates", func(c *fiber.Ctx) error {
		workflowTemplates := argoAdapter.GetWorkflowTemplates()
		var templateGroups []workflow_2.TemplateGroup
		for _, argoTemplate := range workflowTemplates.Items {
			templateGroups = append(templateGroups, *workflow_2.LoadTemplateGroupFromArgo(&argoTemplate))
		}
		return c.JSON(templateGroups)
		// return c.JSON(workflowTemplates)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		body := c.Request().Body()
		var reqBody argo_adapter.WorkflowTemplateSubmitBody
		err := json.Unmarshal(body, &reqBody)
		if err != nil {
			return err
		}
		workflow, err := argoAdapter.SubmitWorkflowTemplate(reqBody)
		if err != nil {
			return err
		}
		return c.JSON(workflow)
	})

	// app.Get("/workflow-events", func(c *fiber.Ctx) error {
	// 	var message = make(chan argo_adapter.WorkflowEvent)
	// 	go argoAdapter.ListenWorkflowEvent("caro-template-j25dz", "argo", message)

	// 	argoWfEvent := <-message

	// 	argoWf := argoWfEvent.Result.Object
	// 	blueprint := workflow.NewArgoBlueprint()
	// 	blueprint.Load(argoWf)
	// 	fmt.Println(blueprint)

	// 	return c.JSON(blueprint)
	// })

	testTemplateModel()
	// testWorkflowModel()
	app.Listen(":3000")

}

func testTemplateModel() {
	// Read Worflow template data from json file
	json_bytes, err := ioutil.ReadFile("workflow-template.json")
	if err != nil {
		fmt.Println(err)
	}

	var argoWfTemplate argo_adapter.WorkflowTemplate
	err = json.Unmarshal(json_bytes, &argoWfTemplate)
	if err != nil {
		fmt.Println(err)
	}

	templateGroup := workflow_2.LoadTemplateGroupFromArgo(&argoWfTemplate)

	str, _ := json.MarshalIndent(templateGroup, "", "\t")
	fmt.Println(string(str))

	// I had templates, which is used for making a blueprint

	blueprint := workflow_2.NewArgoBlueprint()

	blueprint.AddNode(templateGroup.Templates[0].Id, &templateGroup.Templates[0])
	blueprint.AddNode(templateGroup.Templates[1].Id, &templateGroup.Templates[1])
	blueprint.AddEdge(templateGroup.Templates[0].Id, templateGroup.Templates[1].Id)

	fmt.Println(blueprint)

	// Read Worflow data from json file
	json_bytes, err = ioutil.ReadFile("workflow.json")
	if err != nil {
		fmt.Println(err)
	}

	var argoWf argo_adapter.Workflow
	err = json.Unmarshal(json_bytes, &argoWf)
	if err != nil {
		fmt.Println(err)
	}

	loadedBlueprint := workflow_2.LoadBlueprint(argoWf)

	fmt.Println("loadedBlueprint", loadedBlueprint)
	sampleNode := loadedBlueprint.GetNodes()[0]
	fmt.Println(sampleNode)
	fmt.Println(loadedBlueprint.GetNodes()[1])

	// Sample blueprint ???

	// Blueprints use for create a workflow, We persis the blueprint
}

// func testWorkflowModel() {
// 	// Read Worflow data from json file
// 	json_bytes, err := ioutil.ReadFile("workflow.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var argoWfSample argo_adapter.Workflow
// 	err = json.Unmarshal(json_bytes, &argoWfSample)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Create workflow from argo workflow
// 	workflow := workflow_2.LoadWorkflow(argoWfSample)

// 	fmt.Println(workflow.ArgoBlueprint.String())

// 	// Print out a sample node
// 	nodes := workflow.ArgoBlueprint.GetNodes()
// 	sampleNode := nodes[0]
// 	fmt.Println(sampleNode)

// 	// wf := workflow.Submit()
// 	workflow.Watch()
// }
