import matplotlib.pyplot as plt
import numpy as np
import csv

def main():

    data_dict = {}
    selected_tool = "go-wrk"
    row_name = "Total requests made"

    with open("results.csv") as csv_file:
        csv_reader = csv.DictReader(csv_file)

        for row in csv_reader:
            if row["Tool"] == selected_tool:

                value = float(row[row_name])
                if value == 0:
                    continue

                data_dict[row["Technology"]] = value

    #sort dict by values
    data_dict = dict(sorted(data_dict.items(), key=lambda item: item[1]))

    #begin plot
    plt.style.use('_mpl-gallery')

    fig, ax = plt.subplots()

    #add bars
    for technology, avg_latency in data_dict.items():
        bar = ax.bar(technology, avg_latency, 1)

        ax.bar_label(bar, label_type="edge", labels=[technology])

    plt.title("%s measured with %s" % (row_name, selected_tool))
    plt.show()

if __name__ == "__main__":
    main()