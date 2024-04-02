String fileContents = new File('/home/pipelines/firstjob.Jenkinsfile').text
pipelineJob("firstjob") {
    parameters {
        stringParam('name', "", 'name of the person')
    }
    definition {
        cps {
            script(fileContents)
            sandbox()
        }
    }
}