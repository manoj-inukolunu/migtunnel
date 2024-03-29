$(document).ready(function () {
  $("#tunnels-list a").map(function (elem, d, d1) {
    $(this).removeClass();
    $(this).addClass(
        "block py-2.5 px-4 text-white rounded transition duration-200 hover:bg-blue-700 hover:text-white");
  });
  $('#content-container').hide();
  var table = $('#example').DataTable({
    "select": true,
    "processing": true,
    "serverSide": true,
    "ajax": "/rest/data/history",
    "initComplete": function () {
      $('.dataTables_filter input').unbind();
      $('.dataTables_filter input').bind('keyup', function (e) {
        var code = e.keyCode || e.which;
        if (code == 13) {
          table.search(this.value).draw();
        }
      });
    },
    "columns": [
      {"data": "requestId", "width": "85px"},
      {"data": "requestTime", "width": "150px"},
      {
        "data": "line",
        "width": "200px",
        "render": function (data, type, row) {
          return data;
        }
      }
    ]
  });

  $('#example tbody').on('click', 'tr', function () {
    $(this).toggleClass('selected');
  });

  function removeAllClasses() {
    $('#nav-pills-id li a').each(function () {
      $(this).removeClass('active');
    });
  }

  var requestId;
  table.off('click', 'tbody tr').on('click', 'tbody tr', function () {
    $('#content-container').unbind();
    requestId = table.row(this).data().requestId;
    $.get("/request/" + requestId).done(function (content) {
      $('#content').text(content);
    });
    $('#content-container').show();

  });
  $('#responsePane').click(function () {
    removeAllClasses();
    $('#responsePane').addClass('active');
    $.get("/response/" + requestId).done(function (content) {
      $('#content').text(content);
    });
  });
  $('#requestPane').click(function () {
    removeAllClasses();
    $('#requestPane').addClass('active');
    $.get("/request/" + requestId).done(function (content) {
      $('#content').text(content);
    });
  });

  $('#replay').click(function () {
    removeAllClasses();
    $('#replay').addClass('active');
    $.get("/replay/" + requestId).done(function (content) {
      $('#content').text(content);
    });
  });

  $('#delete').click(function () {
    removeAllClasses();
    $('#delete').addClass('active');
    $.get("/delete/" + requestId).done(function (content) {
      location.reload();
    });
  });
});




