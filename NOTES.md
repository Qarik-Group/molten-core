```
storage:
  files:
  - filesystem: "root"
    path: "/opt/bin/mc"
    mode: 0755
    contents:
      source: "https://example.com/metadata-script.sh"
      verification:
        hash: sha512-...  # <type>-<value>

systemd:
  units:
  - name: mc.service
    enable: true
    contents: |
      [Unit]
      Description=MoltenCore - configures cluster for usage with BUCC
      Requires=flannel.service coreos-metadata.service

      [Service]
      Type=oneshot
      ExecStart=/opt/bin/mc init
      RemainAfterExit=true
      StandardOutput=journal

      [Install]
      WantedBy=multi-user.target docker.service
```
