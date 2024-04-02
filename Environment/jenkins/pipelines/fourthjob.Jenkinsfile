pipeline {
    agent any

    stages {
        stage('Run other first jobs in parallel') {
            steps {
                parallel (
                    linux: {
                        build job: 'firstjob'
                    },
                    mac: {
                        build job: 'firstjob'
                    },
                    windows: {
                        build job: 'firstjob'
                    },
                    failFast: false)
            }
        }
        stage('Run third job') {
            steps {
                build job: 'thirdjob', 
                    parameters: [
                        string(name: 'name', value: String.valueOf(BUILD_NUMBER))
                    ]
            }
        }
    }
}