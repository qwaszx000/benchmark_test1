import os
import os.path
import re
import csv
import pandas

def process_wrk_report(file):
    result = [0] * 6

    #skip 3 lines
    for _ in range(3):
        file.readline()

    line_data_list = []
    for data_cell in file.readline().split(" "):
        stripped = data_cell.strip()
        if stripped:
            line_data_list.append(stripped)
    
    avg_latency_str = line_data_list[1]
    max_latency_str = line_data_list[3]

    #skip 6 lines
    for _ in range(6):
        file.readline()

    total_requests_sent = 0
    for data_cell in file.readline().split(" "):
        stripped = data_cell.strip()
        if stripped:
            total_requests_sent = int(stripped)
            break

    line = file.readline().strip()

    avg_rps = 0.0
    errors_count = 0

    def _parse_avg_rps_line(_line):
        for data_cell in reversed(_line.split(" ")):
            stripped = data_cell.strip()
            if stripped:
                return float(stripped)

    #Errors line and/or avg rps line
    is_avg_parsed = False
    if line.startswith("Socket"):
        for data_cell in line.split(" "):
            if data_cell.isdigit():
                errors_count + int(data_cell)
    else:
        avg_rps = _parse_avg_rps_line(line)
        is_avg_parsed = True

    if not is_avg_parsed:
        avg_rps = _parse_avg_rps_line(file.readline().strip())

    result[0] = total_requests_sent
    result[1] = avg_rps

    result[2] = 0 #no data
    result[3] = pandas.Timedelta(avg_latency_str).total_seconds()
    result[4] = pandas.Timedelta(max_latency_str).total_seconds()

    result[5] = (total_requests_sent - errors_count) / total_requests_sent

    return result

def process_go_wrk_report(file):
    result = [0] * 6

    #skip 2 lines
    for _ in range(2):
        file.readline()

    line = file.readline()
    if line.startswith("Error"):
        return result
    
    total_requests_sent = 0
    for data_cell in line.split(" "):
        stripped = data_cell.strip()
        if stripped and stripped.isdigit():
            total_requests_sent = int(stripped)
            break

    #skip 2 lines
    for _ in range(2):
        file.readline()

    avg_rps = float(file.readline().strip().split("\t")[-1])

    #skip 1 line
    file.readline()

    min_latency_str = file.readline().strip().split("\t")[-1]
    avg_latency_str = file.readline().strip().split("\t")[-1]
    max_latency_str = file.readline().strip().split("\t")[-1]

    errors_count = int(file.readline().strip().split("\t")[-1])

    result[0] = total_requests_sent
    result[1] = avg_rps

    result[2] = pandas.Timedelta(min_latency_str).total_seconds()
    result[3] = pandas.Timedelta(avg_latency_str).total_seconds()
    result[4] = pandas.Timedelta(max_latency_str).total_seconds()

    result[5] = (total_requests_sent - errors_count) / total_requests_sent

    return result

def process_my_util_report(file):
    result = [0] * 6

    while True:
        line = file.readline()
        if line.endswith("Done\n"):
            break
    
    #total_requests_sent = int(file.readline().strip().split(" ")[-1])
    #skip 1 line
    file.readline()

    avg_rps = float(file.readline().strip().split(" ")[-1])
    success_rate = float(file.readline().strip().split(" ")[-1][:-1]) / 100

    #skip 1 line
    file.readline()

    min_latency_str = file.readline().strip().split(" ")[-1]
    max_latency_str = file.readline().strip().split(" ")[-1]
    avg_latency_str = file.readline().strip().split(" ")[-1]

    total_requests_sent = int(file.readline().strip().split(" ")[-1])

    result[0] = total_requests_sent
    result[1] = avg_rps

    result[2] = pandas.Timedelta(min_latency_str).total_seconds()
    result[3] = pandas.Timedelta(avg_latency_str).total_seconds()
    result[4] = pandas.Timedelta(max_latency_str).total_seconds()

    result[5] = success_rate

    return result

def process_docker_stats_file(fd):
    result = []

    for line in fd:
        #Filter header lines, we don't need them
        if ord(line[0]) == 0x1b:
            continue

        line_cpu_usage = 0.0
        line_mem_usage = 0

        data_list = line.split(" ")
        for cell in data_list:
            if len(cell) == 0:
                continue

            if cell[-1] == '%':
                cpu_usage = float(cell[:-1])
                line_cpu_usage = cpu_usage
                continue

            if cell[-1] == 'B':
                mem_usage = float(cell[:-3])

                #Handle human-readable format
                if cell[-3:] == "MiB":
                    mem_usage = mem_usage * 1024 * 1024
                elif cell[-3:] == "KiB":
                    mem_usage = mem_usage * 1024

                line_mem_usage = int(mem_usage)

                #We don't need other info
                break
        
        result.append([line_mem_usage, line_cpu_usage])

    return result

def process_file(file_path):
    filename = os.path.basename(file_path)
    technology, test_tool = filename[:-4].split("__")

    max_mem_usage = 0
    max_cpu_usage = 0.0
    for entity in os.scandir("../docker_stats"):
        if entity.is_file() and f"{technology}.csv" == entity.name:

            with open(entity) as top_report_file:
                reader = csv.DictReader(top_report_file)
                
                for data in reader:
                    mem_usage = int(data["Mem usage"])
                    cpu_usage = float(data["CPU%"])
                    if mem_usage > max_mem_usage:
                        max_mem_usage = mem_usage

                    if cpu_usage > max_cpu_usage:
                        max_cpu_usage = cpu_usage
                
            break
    
    data = [technology, test_tool]
    with open(file_path) as report_file:
        if test_tool == "wrk":
            data += process_wrk_report(report_file)
        elif test_tool == "go-wrk":
            data += process_go_wrk_report(report_file)
        elif test_tool == "my_util":
            data += process_my_util_report(report_file)
        else:
            raise ValueError("test_tool is not supported")
    
    return data + [max_mem_usage, max_cpu_usage]

def translate_docker_stats2csv():
    for entity in os.scandir("../docker_stats"):
        if entity.is_file():
            print(entity.name)

            header = ["Mem usage", "CPU%"]
            data = None
            with open(entity.path) as fd:
                data = process_docker_stats_file(fd)
            
            result_csv_path = os.path.join("../docker_stats", entity.name + ".csv")
            with open(result_csv_path, "w") as fd:
                writer = csv.writer(fd)
                writer.writerow(header)
                writer.writerows(data)

def translate_resuls2csv():
    with open("results.csv", "w") as csv_file:
        results_writer = csv.writer(csv_file)

        header_row = [
            "Technology",
            "Tool",

            "Total requests made",
            "AVG RPS",

            "Min latency",
            "AVG latency",
            "Max latency",

            "Success rate",

            "Max mem usage",
            "Max CPU usage"
        ]
        results_writer.writerow(header_row)

        for entity in os.scandir("../results"):
            if entity.is_file():
                print(entity.name)

                csv_row_data = process_file(entity.path)
                results_writer.writerow(csv_row_data)

def main():
    translate_docker_stats2csv()
    translate_resuls2csv()
    
if __name__ == "__main__":
    main()