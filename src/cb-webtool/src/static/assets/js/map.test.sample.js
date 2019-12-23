
            $(document).ready(function()
            {
                //$('#js-page-content').smartPanel(); 

                //
                //
                var dataSetPie = [
                {
                    label: "Asia",
                    data: 4119630000,
                    color: color.primary._500
                },
                {
                    label: "Latin America",
                    data: 590950000,
                    color: color.info._500
                },
                {
                    label: "Africa",
                    data: 1012960000,
                    color: color.warning._500
                },
                {
                    label: "Oceania",
                    data: 95100000,
                    color: color.danger._500
                },
                {
                    label: "Europe",
                    data: 727080000,
                    color: color.success._500
                },
                {
                    label: "North America",
                    data: 344120000,
                    color: color.fusion._400
                }];


                $.plot($("#flotPie"), dataSetPie,
                {
                    series:
                    {
                        pie:
                        {
                            innerRadius: 0.5,
                            show: true,
                            radius: 1,
                            label:
                            {
                                show: true,
                                radius: 2 / 3,
                                threshold: 0.1
                            }
                        }
                    },
                    legend:
                    {
                        show: false
                    }
                });


                $.plot('#flotBar1', [
                {
                    data: [
                        [1, 0],
                        [2, 0],
                        [3, 0],
                        [4, 1],
                        [5, 3],
                        [6, 3],
                        [7, 10],
                        [8, 11],
                        [9, 10],
                        [10, 9],
                        [11, 12],
                        [12, 8],
                        [13, 10],
                        [14, 6],
                        [15, 3]
                    ],
                    bars:
                    {
                        show: true,
                        lineWidth: 0,
                        fillColor: color.fusion._50,
                        barWidth: .3,
                        order: 'left'
                    }
                },
                {
                    data: [
                        [1, 0],
                        [2, 0],
                        [3, 1],
                        [4, 2],
                        [5, 2],
                        [6, 5],
                        [7, 8],
                        [8, 12],
                        [9, 10],
                        [10, 11],
                        [11, 3]
                    ],
                    bars:
                    {
                        show: true,
                        lineWidth: 0,
                        fillColor: color.success._500,
                        barWidth: .3,
                        align: 'right'
                    }
                }],
                {
                    grid:
                    {
                        borderWidth: 0,
                    },
                    yaxis:
                    {
                        min: 0,
                        max: 15,
                        tickColor: '#F0F0F0',
                        ticks: [
                            [0, ''],
                            [5, '$5000'],
                            [10, '$25000'],
                            [15, '$45000']
                        ],
                        font:
                        {
                            color: '#444',
                            size: 10
                        }
                    },
                    xaxis:
                    {
                        mode: 'categories',
                        tickColor: '#F0F0F0',
                        ticks: [
                            [0, '3am'],
                            [1, '4am'],
                            [2, '5am'],
                            [3, '6am'],
                            [4, '7am'],
                            [5, '8am'],
                            [6, '9am'],
                            [7, '10am'],
                            [8, '11am'],
                            [9, '12nn'],
                            [10, '1pm'],
                            [11, '2pm'],
                            [12, '3pm'],
                            [13, '4pm'],
                            [14, '5pm']
                        ],
                        font:
                        {
                            color: '#999',
                            size: 9
                        }
                    }
                });


                /*
                 * VECTOR MAP
                 */

                var data_array = {
                    "af": "16.63",
                    "al": "0",
                    "dz": "158.97",
                    "ao": "85.81",
                    "ag": "1.1",
                    "ar": "351.02",
                    "am": "8.83",
                    "au": "1219.72",
                    "at": "366.26",
                    "az": "52.17",
                    "bs": "7.54",
                    "bh": "21.73",
                    "bd": "105.4",
                    "bb": "3.96",
                    "by": "52.89",
                    "be": "461.33",
                    "bz": "1.43",
                    "bj": "6.49",
                    "bt": "1.4",
                    "bo": "19.18",
                    "ba": "16.2",
                    "bw": "12.5",
                    "br": "2023.53",
                    "bn": "11.96",
                    "bg": "44.84",
                    "bf": "8.67",
                    "bi": "1.47",
                    "kh": "11.36",
                    "cm": "21.88",
                    "ca": "1563.66",
                    "cv": "1.57",
                    "cf": "2.11",
                    "td": "7.59",
                    "cl": "199.18",
                    "cn": "5745.13",
                    "co": "283.11",
                    "km": "0.56",
                    "cd": "12.6",
                    "cg": "11.88",
                    "cr": "35.02",
                    "ci": "22.38",
                    "hr": "59.92",
                    "cy": "22.75",
                    "cz": "195.23",
                    "dk": "304.56",
                    "dj": "1.14",
                    "dm": "0.38",
                    "do": "50.87",
                    "ec": "61.49",
                    "eg": "216.83",
                    "sv": "21.8",
                    "gq": "14.55",
                    "er": "2.25",
                    "ee": "19.22",
                    "et": "30.94",
                    "fj": "3.15",
                    "fi": "231.98",
                    "fr": "2555.44",
                    "ga": "12.56",
                    "gm": "1.04",
                    "ge": "11.23",
                    "de": "3305.9",
                    "gh": "18.06",
                    "gr": "305.01",
                    "gd": "0.65",
                    "gt": "40.77",
                    "gn": "4.34",
                    "gw": "0.83",
                    "gy": "2.2",
                    "ht": "6.5",
                    "hn": "15.34",
                    "hk": "226.49",
                    "hu": "132.28",
                    "is": "0",
                    "in": "1430.02",
                    "id": "695.06",
                    "ir": "337.9",
                    "iq": "84.14",
                    "ie": "204.14",
                    "il": "201.25",
                    "it": "2036.69",
                    "jm": "13.74",
                    "jp": "5390.9",
                    "jo": "27.13",
                    "kz": "129.76",
                    "ke": "32.42",
                    "ki": "0.15",
                    "kw": "117.32",
                    "kg": "4.44",
                    "la": "6.34",
                    "lv": "23.39",
                    "lb": "39.15",
                    "ls": "1.8",
                    "lr": "0.98",
                    "lt": "35.73",
                    "lu": "52.43",
                    "mk": "9.58",
                    "mg": "8.33",
                    "mw": "5.04",
                    "my": "218.95",
                    "mv": "1.43",
                    "ml": "9.08",
                    "mt": "7.8",
                    "mr": "3.49",
                    "mu": "9.43",
                    "mx": "1004.04",
                    "md": "5.36",
                    "rw": "5.69",
                    "ws": "0.55",
                    "st": "0.19",
                    "sa": "434.44",
                    "sn": "12.66",
                    "rs": "38.92",
                    "sc": "0.92",
                    "sl": "1.9",
                    "sg": "217.38",
                    "sk": "86.26",
                    "si": "46.44",
                    "sb": "0.67",
                    "za": "354.41",
                    "es": "1374.78",
                    "lk": "48.24",
                    "kn": "0.56",
                    "lc": "1",
                    "vc": "0.58",
                    "sd": "65.93",
                    "sr": "3.3",
                    "sz": "3.17",
                    "se": "444.59",
                    "ch": "522.44",
                    "sy": "59.63",
                    "tw": "426.98",
                    "tj": "5.58",
                    "tz": "22.43",
                    "th": "312.61",
                    "tl": "0.62",
                    "tg": "3.07",
                    "to": "0.3",
                    "tt": "21.2",
                    "tn": "43.86",
                    "tr": "729.05",
                    "tm": "0",
                    "ug": "17.12",
                    "ua": "136.56",
                    "ae": "239.65",
                    "gb": "2258.57",
                    "us": "14624.18",
                    "uy": "40.71",
                    "uz": "37.72",
                    "vu": "0.72",
                    "ve": "285.21",
                    "vn": "101.99",
                    "ye": "30.02",
                    "zm": "15.69",
                    "zw": "0"
                };

                $('#vector-map').vectorMap(
                {
                    map: 'world_en',
                    backgroundColor: 'transparent',
                    color: color.warning._50,
                    borderOpacity: 0.5,
                    borderWidth: 1,
                    hoverColor: color.success._300,
                    hoverOpacity: null,
                    selectedColor: color.success._500,
                    selectedRegions: ['US'],
                    enableZoom: true,
                    showTooltip: true,
                    scaleColors: [color.primary._400, color.primary._50],
                    values: data_array,
                    normalizeFunction: 'polynomial',
                    onRegionClick: function(element, code, region)
                    {
                        /*var message = 'You clicked "'
						+ region
						+ '" which has the code: '
						+ code.toLowerCase();
			 
					console.log(message);*/

                        var randomNumber = Math.floor(Math.random() * 10000000);
                        var arrow;

                        if (Math.random() >= 0.5 == true)
                        {
                            arrow = '<div class="ml-2 d-inline-flex"><i class="fal fa-caret-up text-success fs-xs"></i></div>'
                        }
                        else
                        {
                            arrow = '<div class="ml-2 d-inline-flex"><i class="fal fa-caret-down text-danger fs-xs"></i></div>'
                        }

                        $('.js-jqvmap-flag').attr('src', 'https://lipis.github.io/flag-icon-css/flags/4x3/' + code.toLowerCase() + '.svg');
                        $('.js-jqvmap-country').html(region + ' - ' + '$' + randomNumber.toFixed(2).replace(/(\d)(?=(\d\d\d)+(?!\d))/g, "$1,") + arrow);
                    }
                });


                /* TAB 1: UPDATING CHART */
                var data = [],
                    totalPoints = 200;
                var getRandomData = function()
                {
                    if (data.length > 0)
                        data = data.slice(1);

                    // do a random walk
                    while (data.length < totalPoints)
                    {
                        var prev = data.length > 0 ? data[data.length - 1] : 50;
                        var y = prev + Math.random() * 10 - 5;
                        if (y < 0)
                            y = 0;
                        if (y > 100)
                            y = 100;
                        data.push(y);
                    }

                    // zip the generated y values with the x values
                    var res = [];
                    for (var i = 0; i < data.length; ++i)
                        res.push([i, data[i]])
                    return res;
                }
                // setup control widget
                var updateInterval = 1500;
                $("#updating-chart").val(updateInterval).change(function()
                {

                    var v = $(this).val();
                    if (v && !isNaN(+v))
                    {
                        updateInterval = +v;
                        $(this).val("" + updateInterval);
                    }

                });
                // setup plot
                var options = {
                    colors: [color.primary._700],
                    series:
                    {
                        lines:
                        {
                            show: true,
                            lineWidth: 0.5,
                            fill: 0.9,
                            fillColor:
                            {
                                colors: [
                                {
                                    opacity: 0.6
                                },
                                {
                                    opacity: 0
                                }]
                            },
                        },

                        shadowSize: 0 // Drawing is faster without shadows
                    },
                    grid:
                    {
                        borderColor: '#F0F0F0',
                        borderWidth: 1,
                        labelMargin: 5
                    },
                    xaxis:
                    {
                        color: '#F0F0F0',
                        font:
                        {
                            size: 10,
                            color: '#999'
                        }
                    },
                    yaxis:
                    {
                        min: 0,
                        max: 100,
                        color: '#F0F0F0',
                        font:
                        {
                            size: 10,
                            color: '#999'
                        }
                    }
                };
                var plot = $.plot($("#updating-chart"), [getRandomData()], options);
                /* live switch */
                $('input[type="checkbox"]#start_interval').click(function()
                {
                    if ($(this).prop('checked'))
                    {
                        $on = true;
                        updateInterval = 1500;
                        update();
                    }
                    else
                    {
                        clearInterval(updateInterval);
                        $on = false;
                    }
                });
                var update = function()
                {
                    if ($on == true)
                    {
                        plot.setData([getRandomData()]);
                        plot.draw();
                        setTimeout(update, updateInterval);

                    }
                    else
                    {
                        clearInterval(updateInterval)
                    }

                }



                /*calendar */
                var todayDate = moment().startOf('day');
                var YM = todayDate.format('YYYY-MM');
                var YESTERDAY = todayDate.clone().subtract(1, 'day').format('YYYY-MM-DD');
                var TODAY = todayDate.format('YYYY-MM-DD');
                var TOMORROW = todayDate.clone().add(1, 'day').format('YYYY-MM-DD');


                var calendarEl = document.getElementById('calendar');

                var calendar = new FullCalendar.Calendar(calendarEl,
                {
                    plugins: ['dayGrid', 'list', 'timeGrid', 'interaction', 'bootstrap'],
                    themeSystem: 'bootstrap',
                    timeZone: 'UTC',
                    dateAlignment: "month", //week, month
                    buttonText:
                    {
                        today: 'today',
                        month: 'month',
                        week: 'week',
                        day: 'day',
                        list: 'list'
                    },
                    eventTimeFormat:
                    {
                        hour: 'numeric',
                        minute: '2-digit',
                        meridiem: 'short'
                    },
                    navLinks: true,
                    header:
                    {
                        left: 'title',
                        center: '',
                        right: 'today prev,next'
                    },
                    footer:
                    {
                        left: 'dayGridMonth,timeGridWeek,timeGridDay,listWeek',
                        center: '',
                        right: ''
                    },
                    editable: true,
                    eventLimit: true, // allow "more" link when too many events
                    events: [
                    {
                        title: 'All Day Event',
                        start: YM + '-01',
                        description: 'This is a test description', //this is currently bugged: https://github.com/fullcalendar/fullcalendar/issues/1795
                        className: "border-warning bg-warning text-dark"
                    },
                    {
                        title: 'Reporting',
                        start: YM + '-14T13:30:00',
                        end: YM + '-14',
                        className: "bg-white border-primary text-primary"
                    },
                    {
                        title: 'Surgery oncall',
                        start: YM + '-02',
                        end: YM + '-03',
                        className: "bg-primary border-primary text-white"
                    },
                    {
                        title: 'NextGen Expo 2019 - Product Release',
                        start: YM + '-03',
                        end: YM + '-05',
                        className: "bg-info border-info text-white"
                    },
                    {
                        title: 'Dinner',
                        start: YM + '-12',
                        end: YM + '-10'
                    },
                    {
                        id: 999,
                        title: 'Repeating Event',
                        start: YM + '-09T16:00:00',
                        className: "bg-danger border-danger text-white"
                    },
                    {
                        id: 1000,
                        title: 'Repeating Event',
                        start: YM + '-16T16:00:00'
                    },
                    {
                        title: 'Conference',
                        start: YESTERDAY,
                        end: TOMORROW,
                        className: "bg-success border-success text-white"
                    },
                    {
                        title: 'Meeting',
                        start: TODAY + 'T10:30:00',
                        end: TODAY + 'T12:30:00',
                        className: "bg-primary text-white border-primary"
                    },
                    {
                        title: 'Lunch',
                        start: TODAY + 'T12:00:00',
                        className: "bg-info border-info"
                    },
                    {
                        title: 'Meeting',
                        start: TODAY + 'T14:30:00',
                        className: "bg-warning text-dark border-warning"
                    },
                    {
                        title: 'Happy Hour',
                        start: TODAY + 'T17:30:00',
                        className: "bg-success border-success text-white"
                    },
                    {
                        title: 'Dinner',
                        start: TODAY + 'T20:00:00',
                        className: "bg-danger border-danger text-white"
                    },
                    {
                        title: 'Birthday Party',
                        start: TOMORROW + 'T07:00:00',
                        className: "bg-primary text-white border-primary text-white"
                    },
                    {
                        title: 'Gotbootstrap Meeting',
                        url: 'http://gotbootstrap.com/',
                        start: YM + '-28',
                        className: "bg-info border-info text-white"
                    }],
                    viewSkeletonRender: function()
                    {
                        $('.fc-toolbar .btn-default').addClass('btn-sm');
                        $('.fc-header-toolbar h2').addClass('fs-md');
                        $('#calendar').addClass('fc-reset-order')
                    },

                });

                calendar.render();
            });

