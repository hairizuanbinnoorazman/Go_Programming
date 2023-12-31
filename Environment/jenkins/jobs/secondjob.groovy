pipelineJob("secondjob") {
    parameters {
        stringParam('name', "", 'name of the person')
    }
    definition {
        cpsScm {
            scm {
                git {
                    remote {
                        url('https://github.com/hairizuanbinnoorazman/Go_Programming')
                    }
                    branch('master')
                }
            }
            scriptPath('Environment/jenkins/pipelines/secondjob.Jenkinsfile')
        }
    }
}
