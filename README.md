## Create own Container

### use linux namespace

1. Namespaces
2. Chroot
3. Cgroups

#### Normal for docker run

docker run --rm -it ubuntu /bin/bash

exit 退出container

##### 使用ps觀察進程發現和本機上的不一樣

- docker container 的隔離機制(原理就是linux 中的namespace)

### Namespace

---

- What you can see
- Created with syscall
  - Unix  Timesharing System
  - Process IDs
  - Mounts
  - Network
  - User IDs
  - InterProcess Comms

### CGroups

---

- What you can use
- Filesystem interface
  - Memory
  - CPU
  - I/O
  - process number

### 程式運行

開始設定container時必須以sudo身份執行

可以使用chroot來限制文件內容
chroot usr
chdir /
ls 無法再看到更上層的目錄

- 使用sleep方法來快速找尋進程

<!-- container -->

sleep 100

<!-- host -->

ps -C sleep
sudo ls -l /proc/[sleep pid]/root
`lrwxrwxrwx 1 root root 0 12月  2 15:17 /proc/60719/root -> /usr`
usr 為我 chroot 的目錄

使用 mount | grep proc 查看mount對象

新增以後發現從主機上mount | grep proc無法查看到container

```go
Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// 希望container具有new ns(namespace) 但不與主機共享
		Unshareflags: syscall.CLONE_NEWNS,
```

docker限制最大記憶體量的方式其實是將數字寫入文件中,這就是它告訴內核它需要限制的最大量

#### 將cgroup降為v1

vim /etc/default/grub
`GRUB_CMDLINE_LINUX="cgroup_enable=memory swapaccount=0 systemd.unified_cgroup_hierarchy=0"`

測試瘋狂創造進程充報記憶體是否還會被cgroup控制組限制->fork bomb
:() { : | : & } ; :

-> 發現永遠不會超過20個進程

查看當前進程
cat pids.current
