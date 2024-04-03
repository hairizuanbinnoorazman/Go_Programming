pipeline {
    agent any

    stages {
        stage('Run second job') {
            steps {
                build job: 'secondjob', 
                    parameters: [
                        string(name: 'name', value: String.valueOf(BUILD_NUMBER))
                    ]
            }
        }
        stage('Other jobs post') {
            steps {
                parallel (
                    linux: {
                        build job: 'thirdjob', 
                            parameters: [
                                string(name: 'name', value: (String.valueOf(BUILD_NUMBER)+'linux'))
                            ], propagate: false, wait: false
                    },
                    mac: {
                        build job: 'thirdjob', 
                            parameters: [
                                string(name: 'name', value: (String.valueOf(BUILD_NUMBER)+'mac'))
                            ], propagate: false, wait: false
                    },
                    windows: {
                        build job: 'thirdjob', 
                            parameters: [
                                string(name: 'name', value: (String.valueOf(BUILD_NUMBER)+'windows'))
                            ], propagate: false, wait: false
                    },
                    failFast: false)
            }
        }
    }
}