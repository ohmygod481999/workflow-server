package main

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"callbot/workflow/auth"
	"callbot/workflow/workflow"
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
		return c.JSON(workflowTemplates)
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

	app.Get("/workflow-events", func(c *fiber.Ctx) error {
		var message = make(chan argo_adapter.WorkflowEvent)
		go argoAdapter.ListenWorkflowEvent("caro-template-j25dz", "argo", message)

		argoWfEvent := <-message

		argoWf := argoWfEvent.Result.Object
		wf := workflow.NewBlueprint()
		wf.ReadFromArgoWorkflow(argoWf)
		fmt.Println(wf)

		return c.JSON(wf)
	})

	testWorkflowModel()
	app.Listen(":3000")

}

func testWorkflowModel() {
	// Read Worflow data from json file
	json_bytes, err := ioutil.ReadFile("workflow.json")
	if err != nil {
		fmt.Println(err)
	}
	var wfSample argo_adapter.Workflow
	err = json.Unmarshal(json_bytes, &wfSample)

	// Create workflow from argo workflow
	blueprint := workflow.NewBlueprint()
	blueprint.ReadFromArgoWorkflow(wfSample)
	fmt.Println(blueprint)

	// Print out a sample node
	nodes := blueprint.GetNodes()
	sampleNode := nodes[0]
	fmt.Println(sampleNode)
}
