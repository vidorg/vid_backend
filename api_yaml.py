import json
import os
import re
import sys

import yaml


def parsePlain(content) -> {}:
    """
    Parse plain @annotion string (2-lier)
    """
    tokens = re.findall(r'// @(.+)', content)
    plain = {}
    for token in tokens:
        sp = re.split(r'[ \t]', token)
        val = ' '.join(sp[1:])
        sp[0] = sp[0].lower()
        if '.' in sp[0]:
            spd = sp[0].split('.')
            if spd[0] not in plain.keys():
                plain[spd[0]] = {}
            plain[spd[0]][spd[1]] = val.strip(' \t')
        else:
            plain[sp[0]] = val.strip(' \t')
    return plain


def parseArray(content, field) -> []:
    """
    Parse plain @annotion array
    """
    tupls = re.findall(r'// @(' + field + '|' + field.lower() +
                       '|' + field.capitalize() + f')(.+)', content)
    return [tupl[1].strip(' \t') for tupl in tupls]


def parseMultiLine(content, field) -> str:
    """
    Parse plain @annotion multiline string
    """
    tupls = re.compile(r'/\* @(' + field.lower() + '|' +
                       field.capitalize() + f')(.+?)\*/', re.DOTALL).findall(content)
    if len(tupls) == 0:
        return ''
    cn = re.sub(r'( \t|\t |\t)', '', tupls[0][1])
    return cn


###


def parseParam(content, need_re=True) -> []:
    """
    Parse @Param name; in; type; required; description; enum; minLength; maxLength
    $name, $in, $type, $req, $desc, [$enum, $minLength, $maxLength]

    Ex: sex formData string false "new sex" enum(male, female, unknown)
    """
    if need_re:
        tokens = re.findall(r'// @Param[ \t]+(.+)', content)
    else:
        tokens = [content]
    params = []
    for token in tokens:
        token = re.split(r'[ \t]', token)
        name, in_type, param_type, is_req, others = token[0], token[1], token[2], token[3], token[4:]

        others = ' '.join(others)
        cmt = re.compile(r'.*"(.*)".*').findall(others)
        cmt = cmt[0].strip(' \t') if len(cmt) != 0 else ''
        enum = re.compile(r'.*enum\((.+?)\).*').findall(others)
        enum = [e.strip(' \t')
                for e in enum[0].split(',')] if len(enum) != 0 else []
        maxLength = re.compile(r'.*maxLength\((.+?)\).*').findall(others)
        maxLength = maxLength[0].strip(' \t') if len(maxLength) != 0 else ''
        minLength = re.compile(r'.*minLength\((.+?)\).*').findall(others)
        minLength = minLength[0].strip(' \t') if len(minLength) != 0 else ''

        pj = {
            "name": name.strip(' \t'),
            "description": cmt.strip(' \t'),
            "in": in_type.strip(' \t'),
            "required": True if is_req.strip(' \t').lower() == 'true' else False,
            "type": param_type.strip(' \t')
        }
        if len(enum) != 0:
            pj['enum'] = enum
        if minLength != '':
            pj['minLength'] = int(minLength)
            pj['description'] += '\n*Minimum length* : ' + minLength
        if maxLength != '':
            pj['maxLength'] = int(maxLength)
            pj['description'] += '\n*Maximum length* : ' + maxLength

        params.append(pj)
    return params


def parseResp(content, demo_json) -> []:
    # /* @ (~~~~~~) */
    tokens = re.compile(
        r'/\* @(success|Success|failed|Failed)(.+?)\*/', re.DOTALL).findall(content)
    c = []
    for token in tokens:
        token = token[1]
        sp = re.split(r'[ \t]', token)
        status, code = sp[0], sp[1]
        j = ' '.join(sp[2:]).strip(' \t')

        # parse demo model
        dms = re.compile(r'@\$(.+?)\$').findall(j)
        for dm in dms:
            if demo_json is not None and dm in demo_json:
                try:
                    j = j.replace(f'@${dm}$', json.dumps(demo_json[dm]))
                except:
                    pass
        # to string
        try:
            j = json.dumps(json.loads(j), indent=4)
        except:
            pass
        c.append({
            'status': status.strip(' \t'),
            'code': code.strip(' \t'),
            'json': j
        })
    return c


def parseErrorCode(content, need_re=True) -> str:
    """
    Parse UserDefined Error Code Message
    """
    if need_re:
        ecs = parseArray(content, 'ErrorCode')
    else:
        ecs = content

    codes, messages = [], []
    for ec in ecs:
        sp = re.split(r'[ \t]', ec)
        codes.append(sp[0].strip(' \t'))
        messages.append(' '.join(sp[1:]).strip(' \t'))
    ec = {}
    for c, m in zip(codes, messages):
        if c in ec:
            ec[c]['description'] += ' / "%s"' % m
        else:
            ec[c] = {
                'description': '"%s"' % m
            }
    return ec


###


def parse_main(content):
    """
    parse main.go swagger information
    """
    plain = parsePlain(content)
    out = {}
    if 'basepath' in plain:
        out['basePath'] = plain.pop('basepath').strip(' \t')
    if 'host' in plain:
        out['host'] = plain.pop('host').strip(' \t')
    if 'swagger' in plain:
        out['swagger'] = plain.pop('swagger').strip(' \t')
    out['info'] = plain
    if 'termsofservice' in plain:
        plain['termsOfService'] = plain.pop('termsofservice')
    return out


def parse_demo_response(path) -> {}:
    """
    parse demo.json model default response
    """
    demo_json = open(path, 'r', encoding='utf-8').read()
    try:
        return json.loads(demo_json)
    except:
        return None


def parse_ctrl(content, out_yaml, auth_param, auth_error, demo_json):
    """
    parse each of the .go file, go through all functions and parse information
    :param auth_param: common auth param (header)
    :param auth_error: common auth error (401...)
    :param demo_json:  common demo response json model
    """
    if 'paths' not in out_yaml:
        out_yaml['paths'] = {}

    contents = content.split('func ')
    for content in contents:
        try:
            plains = parsePlain(content)
            parameters = parseParam(content)
            c = parseResp(content, demo_json)
            accept = parseArray(content, 'accept')
            tags = parseArray(content, 'tag')
            desc = parseMultiLine(content, 'description')
            if desc == '' and 'description' in plains:
                desc = plains['description']

            router, *route_setting = re.split(r'[ \t]', plains['router'])
            method = route_setting[0][1:-1].lower()
            is_auth = len(route_setting) >= 2 and route_setting[1].lower() == '[auth]'
            oid = router.lower().replace('/', '-') + '-' + method
            oid = oid.replace('--', '-')

            responses = {code['code']: {
                'description': '```json\n%s\n```' % code['json']
            } for code in c}
            responses.update(parseErrorCode(content))

            if is_auth and auth_param is not None:
                for p in auth_param:
                    parameters.insert(0, p)
            if is_auth and auth_error is not None:
                responses.update(auth_error)

            yml = {
                'tags': tags,
                'summary': plains['summary'],
                'description': desc,
                'consumes': accept,
                'produces': ['application/json'],
                'operationId': oid,
                'parameters': parameters,
                'responses': responses,
            }
            if is_auth:
                yml['security'] = [{
                    'basicAuth': ''
                }]

            if router not in out_yaml['paths']:
                out_yaml['paths'][router] = {}
            out_yaml['paths'][router][method] = yml
        except:
            continue


def main():
    main_file = sys.argv[1]
    all_files = []
    for root, _, files in os.walk('.'):
        for f in files:
            if f.split('.')[-1] == 'go':
                all_files.append(os.path.join(root, f))

    content = open(main_file, 'r', encoding='utf-8').read()
    print(f"> Parsing {main_file}...")
    out_yaml = parse_main(content)

    # Demo response
    if 'response' in out_yaml['info']:
        demo_path = out_yaml['info'].pop('response')['demopath']
        print(f"> Parsing {demo_path}...")
        demo_json = parse_demo_response(demo_path)
    else:
        demo_json = None

    # Global Auth
    if 'authorization' in out_yaml['info']:
        auth_param = out_yaml['info'].pop('authorization')['param']
        auth_error = parseArray(content, r'authorization\.error')
        auth_param = parseParam(auth_param, need_re=False)
        auth_error = parseErrorCode(auth_error, need_re=False)
    else:
        auth_param, auth_error = None, None

    print(f"> Parsing *.go...")
    for f in all_files:
        content = open(f, 'r', encoding='utf-8').read()
        parse_ctrl(content, out_yaml, auth_param, auth_error, demo_json)

    with open(sys.argv[2], 'w', encoding='utf-8') as f:
        yaml.dump(out_yaml, stream=f, encoding='utf-8', allow_unicode=True)


if __name__ == "__main__":
    main()

# python api_yaml.py main.go ./docs/api.yaml
