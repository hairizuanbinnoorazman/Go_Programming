pipelineJob("thirdjob") {
    description("""
<h1>Third Job</h1>
<p>This is a Third job</p>""")
    parameters {
        stringParam('name', "", 'name of the person')
    }
    triggers {
        parameterizedTimerTrigger {
            parameterizedSpecification('''*/2 * * * * % name=Hola
*/3 * * * * % name=mars''')
        }
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
            scriptPath('Environment/jenkins/pipelines/thirdjob.Jenkinsfile')
        }
    }
}
