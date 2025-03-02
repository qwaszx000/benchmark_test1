import os
import os.path
import time
import subprocess

#rust_faf is not working, so we exclude it for now
exclude_dirs = ["utils", "results", "venv", ".vscode", ".git", "rust_faf", "docker_stats", "imgs"]

local_port = 8081
target_url = "http://127.0.0.1:%d/test_plain" % local_port

results_dir = "results"
docker_stats_dir = "docker_stats"

benchmarking_cmds = {
    "wrk": ["utils/wrk", "-c", "50", "-t", "2", "-d", "10s", "--latency", target_url],
    "go-wrk": ["utils/go-wrk", "-c", "50", "-d", "10", "-cpus", "2", target_url],
    "my_util": ["utils/load_simulator/load_simulator", "-rps", "0", "-workers", "50", "-target", target_url]
}

memory_limit = 256*1024*1024 #256MiB

def build_image(dir_name):
    docker_build_env = {
        "DOCKER_BUILDKIT": "1"
    }
    docker_build_cmd = ["docker", "build", "-q"]

    _cmd = docker_build_cmd + [dir_name] #functional append -> return new list instead of changing initial one
    docker_build_instance = subprocess.Popen(_cmd, bufsize=1, text=True, stdout=subprocess.PIPE, env=docker_build_env)
    image_id_str, _ = docker_build_instance.communicate()
    image_id_str = image_id_str.strip()

    return image_id_str

def get_container_id():
    docker_ps_cmd = ["docker", "ps", "-q", "--filter", "label=benchmark_run=1"]
    proc = subprocess.Popen(docker_ps_cmd, bufsize=1, text=True, stdout=subprocess.PIPE)
    container_id_str, _ = proc.communicate()
    container_id_str = container_id_str.strip()

    return container_id_str

def log_docker_stats_to(fd):
    docker_stats_cmd = ["docker", "stats"]
    proc = subprocess.Popen(docker_stats_cmd, bufsize=1, text=True, stdout=fd, stderr=subprocess.STDOUT)

    return proc

def docker_run(dir_name):
    docker_run_cmd = [
        "docker", "run",
        "--memory", str(memory_limit), "--memory-swap", str(memory_limit*2),
        "--publish", "%d:8080" % local_port,
        "--label", "benchmark_run=1"
    ]

    #Build docker container and get it's id
    image_id_str = build_image(dir_name)
    print(image_id_str)

    #Run docker container, sleep to let it start
    _cmd = docker_run_cmd + [image_id_str]
    print(_cmd)
    docker_run_instance = subprocess.Popen(_cmd)
    time.sleep(10)

    #Get container id to stop
    container_id_str = get_container_id()
    print(container_id_str)

    def stop_function():
        docker_stop_cmd = ["docker", "stop"]
        subprocess.run(docker_stop_cmd + [container_id_str])

    return stop_function

def docker_compose_run(dir_name):

    _cmd = [
        "docker-compose", "up", "-d"
    ]

    print(_cmd)
    subprocess.run(_cmd, cwd=dir_name)

    def stop_function():
        _cmd = ["docker-compose", "down"]
        subprocess.run(_cmd, cwd=dir_name)

    return stop_function


docker_compose_apps = ["php_symfony"]
def run_benchmark(dir_name):
    print(dir_name)

    if dir_name in docker_compose_apps:
        stop_docker_fn = docker_compose_run(dir_name)
    else:
        stop_docker_fn = docker_run(dir_name)

    docker_stats_file_path = os.path.join(docker_stats_dir, dir_name)
    docker_stats_fd = open(docker_stats_file_path, "w")
    stats_proc = log_docker_stats_to(docker_stats_fd)

    #Quiet time to log normal memory and cpu usage
    time.sleep(5)

    #Run stress testing tools with redirecting output to txt files
    for tool_name, _tool_cmd in benchmarking_cmds.items():
        result_file_name = "%s__%s.txt" % (dir_name, tool_name)
        result_file_path = os.path.join(results_dir, result_file_name)

        with open(result_file_path, "w") as f:
            subprocess.run(_tool_cmd, bufsize=1, text=True, stdout=f, stderr=subprocess.STDOUT)

    stop_docker_fn()

    time.sleep(5)
    stats_proc.terminate()
    try:
        stats_proc.wait(10)
    except:
        stats_proc.kill()
    docker_stats_fd.close()

def main():
    for entity in os.scandir("./"):
        if entity.is_dir() and entity.name not in exclude_dirs:
            run_benchmark(entity.name)

if __name__ == "__main__":
    main()