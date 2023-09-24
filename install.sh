#os=$(uname)
##cpu_info=$(cat /proc/cpuinfo | grep "model name")
#
#if [[ "$os" == "Darwin" ]]; then
#  cpu_info=$(sysctl -n machdep.cpu.brand_string)
#  echo $cpu_info
#elif [[ "$os" == "Linux" ]]; then
#  echo "Операционная система: Linux"
#elif [[ "$os" == "CYGWIN"* || "$os" == "MINGW"* || "$os" == "MSYS"* ]]; then
#  echo "Операционная система: Windows"
#else
#  echo "Не удалось определить вид операционной системы"
#fi