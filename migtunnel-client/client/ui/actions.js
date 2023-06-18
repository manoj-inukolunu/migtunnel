const actionsDiv = {
    props: ['request'],
    reqeustId: '',
    editor: undefined,
    activeClass: undefined,
    inactiveClass: undefined,
    created() {
        $('#editor').hide();
        $('#request-header').hide();
        emitter.on('getRequestData', function (data) {
            requestId = data.requestId;
            $.get("/request/" + requestId + "/requestData").done(function (content) {
                $('#editor').show();
                $('#request-header').show();
                $('#editor-div').show();
                editor.setValue(content);
                editor.clearSelection();
                editor.session.insert({row: 0, column: 0}, "\r\n");
            });
        })
    },
    mounted() {
        activeClass = 'active inline-block rounded-t-lg bg-gray-100 p-4 text-blue-600 dark:bg-gray-800 dark:text-blue-500';
        inactiveClass = 'inline-block rounded-t-lg p-4 hover:bg-gray-50 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300';
        editor = ace.edit("editor");
        editor.setTheme("ace/theme/monokai");
        // editor.session.setMode("ace/mode/java");
        editor.setReadOnly(true);
        editor.session.setUseWrapMode(true);
        $('#editor').hide();
        $('#request-header').hide();
        $('#editor-div').hide();
    },
    methods: {
        replay() {
            $.get("/replay/" + requestId).done(function (content) {
                editor.setValue(content);
            });
        },
        showrequest() {
            $.get("/request/" + requestId + "/requestData").done(function (content) {
                $('#editor').show();
                $('#request-header').show();
                $('#editor-div').show();
                $('#response-button').removeClass();
                $('#request-button').removeClass();
                $('#response-button').addClass(inactiveClass);
                $('#request-button').addClass(activeClass);
                editor.setValue(content);
                editor.clearSelection();
                editor.session.insert({row: 0, column: 0}, "\r\n");
            });
        },
        showresponse() {
            $.get("/request/" + requestId + "/responseData").done(function (content) {
                $('#editor').show();
                $('#request-header').show();
                $('#editor-div').show();
                $('#response-button').removeClass();
                $('#request-button').removeClass();
                $('#response-button').addClass(activeClass);
                $('#request-button').addClass(inactiveClass);
                editor.setValue(content);
                editor.clearSelection();
                editor.session.insert({row: 0, column: 0}, "\r\n");
            });
        },
        delete() {
            $.get("/delete/" + requestId).done(function (content) {
                location.reload();
            });
        }
    },
    template: `
    <div class="flex-1">
  <header class="my-8 flex flex-col" id="request-header">
    <div class="flex items-center justify-between px-6">
      <button @click="replay" class="mx-2 w-40 rounded border border-gray-400 bg-green-500 py-1 font-semibold text-white shadow hover:bg-green-600">Replay</button>
    </div>
  </header>
  <div class="m-6" id="editor-div">
    <ul class="flex flex-wrap border-b border-gray-200 text-center text-sm font-medium text-gray-500 dark:border-gray-700 dark:text-gray-400">
      <li class="mr-2">
        <a href="#" @click="showrequest" aria-current="page" id="request-button" 
        class="active inline-block rounded-t-lg bg-gray-100 p-4 text-blue-600 dark:bg-gray-800 dark:text-blue-500">Request</a>
      </li>
      <li class="mr-2">
        <a href="#" @click="showresponse" id="response-button" 
        class="inline-block rounded-t-lg p-4 hover:bg-gray-50 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300">Response</a>
      </li>
    </ul>
    <div id="editor" class="mx-5 h-3/4 text-base">Testing</div>
  </div>
</div>
  `
}