pipelineJob("secondjob") {
    description("""
<h1>Second Job</h1>
<p>This is a second job</p>""")
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
