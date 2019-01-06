node {
  def shielderImage

  stage('Git Pull') {
    git url: 'https://github.com/joshchu00/finance-go-shielder.git', branch: 'develop'
  }
  stage('Go Build') {
    sh "${tool name: 'go-1.11', type: 'go'}/bin/go build -a -o main"
  }
  stage('Docker Build') {
    docker.withTool('docker-latest') {
      shielderImage = docker.build('docker.io/joshchu00/finance-go-shielder')
    }
  }
  stage('Docker Push') {
    docker.withTool('docker-latest') {
      docker.withRegistry('', 'DockerHub') {
        shielderImage.push()
      }
    }
  }
}
