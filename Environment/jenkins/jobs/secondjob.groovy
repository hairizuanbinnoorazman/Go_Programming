pipelineJob("secondjob") {
    parameters {
        stringParam('name', "", 'name of the person')
    }
    definition {
        cpsScm {
            scm {
                git('git@github.com:hairizuanbinnoorazman/Go_Programming.git')
            }
            scriptPath('Environment/jenkins/pipelines/secondjob.Jenkinsfile')
        }
    }
}
