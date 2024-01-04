FROM jenkins/agent
USER root
RUN mkdir -p /home/jenkins/.ssh && chown jenkins:jenkins /home/jenkins/.ssh
RUN apt update && apt install -y openssh-server
RUN ssh-keygen -A && service ssh --full-restart
CMD ["/usr/sbin/sshd", "-D"]