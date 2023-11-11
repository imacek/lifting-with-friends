/* globals Chart:false */

(() => {
  'use strict'

  let data = null
  let selectedMetric = null
  let selectedGrouping = null

  const metricSelectBtnElems = document.getElementsByClassName('metric-select-btn')
  const groupingSelectBtnElems = document.getElementsByClassName('grouping-select-btn')

  function handleSelectButtonClicked(btnElems, clickedBtnElem) {
    const selectedValue = clickedBtnElem.getAttribute('value')

    Array.from(btnElems).forEach(element => {
      if (element.getAttribute('value') === selectedValue) {
        element.classList.add('btn-primary')
        element.classList.remove('btn-outline-secondary')
      } else {
        element.classList.add('btn-outline-secondary')
        element.classList.remove('btn-primary')
      }
    })

    return selectedValue
  }

  selectedMetric = handleSelectButtonClicked(metricSelectBtnElems, Array.from(metricSelectBtnElems)[0])
  selectedGrouping = handleSelectButtonClicked(groupingSelectBtnElems, Array.from(groupingSelectBtnElems)[1])

  Array.from(metricSelectBtnElems).forEach(element => {
    element.onclick = event => {
      selectedMetric = handleSelectButtonClicked(metricSelectBtnElems, event.target)
      reloadCharts()
    }
  })
  Array.from(groupingSelectBtnElems).forEach(element => {
    element.onclick = event => {
      selectedGrouping = handleSelectButtonClicked(groupingSelectBtnElems, event.target)
      reloadCharts()
    }
  })

  let chartInstances = []

  function showChart(canvasIndex, exercise, metric, selectedGrouping, data) {
    const canvasElem = document.getElementsByClassName('exercise-chart')[canvasIndex];

    const datasets = Object.entries(data).map(([user, userData]) => {
      return { label: user, data: userData[selectedGrouping][exercise] }
    })

    if (chartInstances[canvasIndex]) {
      chartInstances[canvasIndex].destroy()
    }

    chartInstances[canvasIndex] = new Chart(canvasElem, {
      type: 'line',
      data: {
        datasets: datasets
      },
      options: {
        maintainAspectRatio: false,
        animation: false,
        scales: {
          x: {
            type: 'time',
          }
        },
        parsing: {
          xAxisKey: 'timestamp',
          yAxisKey: metric
        },
        plugins: {
          title: {
            display: true,
            text: exercise,
            color: '#dee2e6',
            font: {
              size: 22
            }
          },
          legend: {
            display: true
          }
        }
      }
    })
  }

  function reloadCharts() {
    if (!data) {
      return
    }

    showChart(0, 'Squat (Barbell)', selectedMetric, selectedGrouping, data)
    showChart(1, 'Deadlift (Barbell)', selectedMetric, selectedGrouping, data)
    showChart(2, 'Bench Press (Barbell)', selectedMetric, selectedGrouping, data)
    showChart(3, 'Overhead Press (Barbell)', selectedMetric, selectedGrouping, data)
    showChart(4, 'Bent Over Row (Barbell)', selectedMetric, selectedGrouping, data)
  }

  async function uploadExerciseCsv()
  {
    let userElem = document.getElementById("user-input")
    let fileElem = document.getElementById("file-input")

    let formData = new FormData();
    formData.append("user", userElem.value);
    formData.append("file", fileElem.files[0]);

    const abortCtrl = new AbortController()
    setTimeout(() => abortCtrl.abort(), 5000);

    try {
      let r = await fetch('/api/upload', { method: "POST", body: formData, signal: abortCtrl.signal });
      console.log('HTTP response code:',r.status);
    } catch(e) {
      console.log('Huston we have problem...:', e);
    }
  }

  //document.getElementById("upload-form").onsubmit = uploadExerciseCsv;


  fetch('/api/data', {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
    },
  })
    .then(response => response.json())
    .then(response => {
      data = response
      reloadCharts()
    })

})()
