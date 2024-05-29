#!/usr/bin/env groovy

def microservices = 'comment, rest, session, user, auth, gateway'.split(', ')

pipeline {
  agent any 

  tools {go "recent go"}

  stages {
    stage('Initialize') {
    steps {
      script {
      for (int i = 0; i < microservices.length; i++) {
        when { changeset pattern: "*${microservices[i]}*", comparator: "REGEXP" } 
          stage("Build Microservice: ${microservices[i]}") {

                script {
                 // sh "sudo cp /home/kayros/backend/config/config.yaml ./config/"
                 // sh "sudo docker build -t resto-${microservices[i]}-service:latest -f ./integration/microservices/${microservices[i]}/Dockerfile ."
  
                  sh "sudo docker-compose -f /home/kayros/backend/integration/prod-compose up --no-deps --build ${microservices[i]}-grpc"
                }
          }
      }
      
      stage('Test') {
            script {
                // sh 'go test ./... -coverprofile=cover.out'
            }
        }

      stage('Code Analysis') {
          script {
            sh 'make easyjs'
            sh '/home/kayros/go/bin/golangci-lint run'
          }
      }

      // for (int i = 0; i < microservices.length; i++) {
      //   stage("Push Microservice: ${microservices[i]}") {
      //           script {
      //             def localImage = "resto-${microservices[i]}-service:latest"
      //             def repositoryName = "kayrosteam/${localImage}"
      //             sh "docker tag ${localImage} ${repositoryName} "
      //             docker.withRegistry("", "dockerhub-credentials") {
      //               def image = docker.image("${repositoryName}");
      //               image.push()
      //             }
      //           }
      //     }
      // } 

      }
    }
    }
  }
}
