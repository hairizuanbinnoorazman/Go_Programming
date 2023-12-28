# Jenkins Setup

This is mostly a jenkins setup in order to gain some experience about how to setup Jenkins

# References

Here are some of the useful references

- https://www.eficode.com/blog/start-jenkins-config-as-code
  - Set the environment variable for Jenkins Configuration as Code to pick up the file with its settings
- https://github.com/jenkinsci/configuration-as-code-plugin/issues/192
  - `matrix_auth` plugin required
- https://github.com/jenkinsci/configuration-as-code-plugin/tree/master/demos/global-matrix-auth
- https://github.com/jenkinsci/configuration-as-code-plugin/tree/master/demos/role-strategy-auth
  - `role_strategy` plugin required
- https://jenkinsci.github.io/job-dsl-plugin/
  - Pretty useful page to understand how to create the initial pipeline jobs