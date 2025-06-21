pipelineJob("fourthjob") {
    description("""
<h1>Fourth Job</h1>
<p>This is a fourth job. It call other jobs</p>""")
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
            scriptPath('Environment/jenkins/pipelines/fourthjob.Jenkinsfile')
        }
    }
}
