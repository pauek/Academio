
Preparación Jutge
-----------------

VMs
---

1. Descargar: 
   - http://ftp.nluug.nl/os/Linux/distr/tinycorelinux/4.x/x86/release/src/kernel/config-3.0.21-tinycore64
   - http://ftp.nluug.nl/os/Linux/distr/tinycorelinux/4.x/x86/release/src/kernel/linux-3.0.21-patched.txz
   - http://ftp.nluug.nl/os/Linux/distr/tinycorelinux/4.x/x86/release/distribution_files/core64.gz

2. Descomprimir el kernel, copiar 'config-3.0.21-tinycore64' a 
   '.config' y marcar las opciones: 
   
     CONFIG_NET_9P=y
     CONFIG_NET_9P_VIRTIO=y  <--- No módulo ('m')!
     CONFIG_9P_FS=y
     CONFIG_9P_FS_POSIX_ACL=y

   Esto implica marcar otras opciones que no estan (con make menuconfig).

3. Compilar el kernel.


[ Ahora se puede botar un Tiny Core Linux con: ]
[   kvm -kernel bzImage -initrd core64.tgz     ]


4. Descomprimir el initrd:

   # mkdir initrd; push initrd
   # zcat ../core64.tgz | sudo cpio -o -H newc -d
   # popd


5. Añadir el driver:

   - Copiar el driver en /usr/bin/driver
   - Editar el /etc/inittab:
     
       + ttyS0::respawn:/usr/bin/driver

[ Ahora se puede conectar con el driver a través  ]
[ del puerto serie: "kvm ... -serial stdio"       ]

6. Compartir ficheros:

   - Crear un directorio para compartir ("grz"?).
   - Añadir la opción -virtfs al lanzar kvm:

       -virtfs local,id=grz,path=grz,security_model=mapped,mount_tag=grz

   - Montar 'grz' en el guest (añadir a /etc/fstab):

       grz  /grz  9p  trans=virtio,version=9p2000.L  0  0


