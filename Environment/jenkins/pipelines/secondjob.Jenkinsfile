pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sleep 10
                echo 'Complete Build step'  
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                sleep 10
                echo 'Complete Testing step'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                sleep 10
                echo 'Complete Deploy step'
            }
        }
    }
}