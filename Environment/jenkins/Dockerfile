FROM jenkins/jenkins:latest
COPY plugins.txt /var/jenkins_home/plugins.txt
RUN /bin/jenkins-plugin-cli -f /var/jenkins_home/plugins.txt
COPY jenkins.yaml /var/jenkins_home/jenkins.yaml
ENV JAVA_OPTS "-Djenkins.install.runSetupWizard=false ${JAVA_OPTS:-}"
ENV CASC_JENKINS_CONFIG=/var/jenkins_home/jenkins.yaml
ENV SSH_PRIVATE_FILE_PATH=/home/jenkins/.ssh/ultimate_ssh_key
RUN git config --global user.email "you@example.com" && \
    git config --global user.name "Your Name"
COPY jobs /home/jobs
COPY pipelines /home/pipelines
COPY --chown=jenkins:jenkins key /home/jenkins/.ssh/ultimate_ssh_key