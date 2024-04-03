pipeline {
    agent any

    stages {
        stage('Printing') {
            steps {
                echo 'Start job...'
                sleep 10
                echo "Name: ${name}"
                script {
                    currentBuild.displayName = "${name}"
                }
            }
        }
    }
}