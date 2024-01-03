FROM jenkins/agent
COPY --chown=jenkins:jenkins key.pub /home/jenkins/.ssh/authorized_keys
USER root
RUN apt update && apt install -y openssh-server
RUN ssh-keygen -A && service ssh --full-restart
CMD ["/usr/sbin/sshd", "-D"]