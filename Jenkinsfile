pipeline{
    agent{
        node{
            label 'linux && go'
        }
    }
    triggers{
        pollSCM("*/20 * * * *")
    }
    options{
        buildDiscarder(logRotator(numToKeepStr: '3'))
    }
    stages {
        stage('build'){
            steps{
                sh 'go build -o main cmd/app/main.go'
            }
        }
        stage('test'){
            environment {
                TEST = credentials('golang-dot-api-test')
            }
            steps{
                sh 'cp -p $TEST $WORKSPACE/test/integration'
                sh 'docker compose -f docker-compose.test.yaml up --build --abort-on-container-exit'
                sh 'docker compose -f docker-compose.test.yaml down --volumes'
                sh 'docker image prune -f'
            }
        }
        stage('deploy'){
            environment {
                DEV = credentials('dev-dot-api')
            }
            steps{
                sh 'mv $DEV .env'
                sh 'docker compose up --build -d'
                sh 'docker image prune -f'
            }
        }
    }
    post{
        success{
            slackSend(color: "good", message: "${env.JOB_NAME} - ${env.BUILD_DISPLAY_NAME} Success after ${currentBuild.durationString.replace(' and counting', '')} (<${env.BUILD_URL}|Open>)")
        }
        failure{
            slackSend(color: "danger", message: "${env.JOB_NAME} - ${env.BUILD_DISPLAY_NAME} Failure after ${currentBuild.durationString.replace(' and counting', '')} (<${env.BUILD_URL}|Open>)")
        }
    }
}